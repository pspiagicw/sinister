package manage

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pspiagicw/sinister/database"
)

// --- pure function unit tests ---

func TestApplyLimit(t *testing.T) {
	input := makeStringEntries(5)

	if got := applyLimit(input, 0); len(got) != 5 {
		t.Errorf("limit=0: expected 5, got %d", len(got))
	}
	if got := applyLimit(input, 3); len(got) != 3 {
		t.Errorf("limit=3: expected 3, got %d", len(got))
	}
	if got := applyLimit(input, 10); len(got) != 5 {
		t.Errorf("limit=10 (over): expected 5, got %d", len(got))
	}
	if got := applyLimit(input, 5); len(got) != 5 {
		t.Errorf("limit=exact: expected 5, got %d", len(got))
	}
}

func TestWithinSinceDays(t *testing.T) {
	recent := time.Now().AddDate(0, 0, -2).UTC().Format(time.RFC3339)
	old := time.Now().AddDate(0, 0, -10).UTC().Format(time.RFC3339)

	if !withinSinceDays(recent, 5) {
		t.Error("recent entry should be within 5 days")
	}
	if withinSinceDays(old, 5) {
		t.Error("old entry should not be within 5 days")
	}
	if !withinSinceDays(old, 0) {
		t.Error("sinceDays=0 should include everything")
	}
	if !withinSinceDays("", 0) {
		t.Error("empty published with sinceDays=0 should be included")
	}
	if withinSinceDays("not-a-date", 5) {
		t.Error("unparseable date with sinceDays>0 should be excluded")
	}
}

// --- integration tests with mock HTTP server ---

func makeFetcher(server *httptest.Server) func(string, int) ([]byte, error) {
	return func(url string, _ int) ([]byte, error) {
		resp, err := http.Get(server.URL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("%s", resp.Status)
		}
		buf := make([]byte, 0, 1024)
		tmp := make([]byte, 512)
		for {
			n, readErr := resp.Body.Read(tmp)
			buf = append(buf, tmp[:n]...)
			if readErr != nil {
				break
			}
		}
		return buf, nil
	}
}

func serveRSS(t *testing.T, xml string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/atom+xml")
		fmt.Fprint(w, xml)
	}))
}

func TestUpdate_InsertsNewEntries(t *testing.T) {
	setupTestDB(t)

	feed := sampleRSSFeed("TestAuthor",
		rssEntry("Video One", "https://www.youtube.com/watch?v=aaa", daysAgo(1))+
			rssEntry("Video Two", "https://www.youtube.com/watch?v=bbb", daysAgo(2)),
	)
	srv := serveRSS(t, feed)
	defer srv.Close()

	Update(UpdateOptions{
		URLs:    []string{srv.URL},
		Fetcher: makeFetcher(srv),
	})

	if database.TotalEntries() != 2 {
		t.Errorf("expected 2 entries, got %d", database.TotalEntries())
	}
	if database.CountUnwatched() != 2 {
		t.Errorf("expected 2 unwatched, got %d", database.CountUnwatched())
	}
}

func TestUpdate_SkipsDuplicates(t *testing.T) {
	setupTestDB(t)

	feed := sampleRSSFeed("TestAuthor",
		rssEntry("Same Video", "https://www.youtube.com/watch?v=dup1", daysAgo(1)),
	)
	srv := serveRSS(t, feed)
	defer srv.Close()

	opts := UpdateOptions{URLs: []string{srv.URL}, Fetcher: makeFetcher(srv)}
	Update(opts)
	Update(opts)

	if database.TotalEntries() != 1 {
		t.Errorf("expected 1 entry after duplicate update, got %d", database.TotalEntries())
	}
}

func TestUpdate_SinceDaysFilter(t *testing.T) {
	setupTestDB(t)

	feed := sampleRSSFeed("TestAuthor",
		rssEntry("Recent", "https://www.youtube.com/watch?v=r1", daysAgo(1))+
			rssEntry("Old", "https://www.youtube.com/watch?v=o1", daysAgo(15)),
	)
	srv := serveRSS(t, feed)
	defer srv.Close()

	Update(UpdateOptions{
		URLs:      []string{srv.URL},
		SinceDays: 7,
		Fetcher:   makeFetcher(srv),
	})

	if database.TotalEntries() != 1 {
		t.Errorf("expected 1 entry (recent only), got %d", database.TotalEntries())
	}
}

func TestUpdate_Limit(t *testing.T) {
	setupTestDB(t)

	feed := sampleRSSFeed("TestAuthor",
		rssEntry("V1", "https://www.youtube.com/watch?v=l1", daysAgo(1))+
			rssEntry("V2", "https://www.youtube.com/watch?v=l2", daysAgo(2))+
			rssEntry("V3", "https://www.youtube.com/watch?v=l3", daysAgo(3)),
	)
	srv := serveRSS(t, feed)
	defer srv.Close()

	Update(UpdateOptions{
		URLs:    []string{srv.URL},
		Limit:   2,
		Fetcher: makeFetcher(srv),
	})

	if database.TotalEntries() != 2 {
		t.Errorf("expected 2 entries with limit=2, got %d", database.TotalEntries())
	}
}

func TestUpdate_DryRun(t *testing.T) {
	setupTestDB(t)

	feed := sampleRSSFeed("TestAuthor",
		rssEntry("Dry Video", "https://www.youtube.com/watch?v=dry1", daysAgo(1)),
	)
	srv := serveRSS(t, feed)
	defer srv.Close()

	Update(UpdateOptions{
		URLs:    []string{srv.URL},
		DryRun:  true,
		Fetcher: makeFetcher(srv),
	})

	if database.TotalEntries() != 0 {
		t.Error("dry run should not insert entries")
	}
}

func TestUpdate_FeedError(t *testing.T) {
	setupTestDB(t)

	good := sampleRSSFeed("GoodAuthor",
		rssEntry("Good Video", "https://www.youtube.com/watch?v=good1", daysAgo(1)),
	)
	goodSrv := serveRSS(t, good)
	defer goodSrv.Close()

	// fetcher that returns error for the bad URL and proxies for the good URL
	fetcher := func(url string, _ int) ([]byte, error) {
		if url == "https://example.com/bad-feed" {
			return nil, fmt.Errorf("404 Not Found")
		}
		resp, err := http.Get(goodSrv.URL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		buf := make([]byte, 0)
		tmp := make([]byte, 512)
		for {
			n, readErr := resp.Body.Read(tmp)
			buf = append(buf, tmp[:n]...)
			if readErr != nil {
				break
			}
		}
		return buf, nil
	}

	Update(UpdateOptions{
		URLs:    []string{"https://example.com/bad-feed", goodSrv.URL},
		Retries: 0,
		Fetcher: fetcher,
	})

	if database.TotalEntries() != 1 {
		t.Errorf("expected 1 entry from good feed, got %d", database.TotalEntries())
	}
}

func TestUpdate_Retry(t *testing.T) {
	setupTestDB(t)

	attempts := 0
	feed := sampleRSSFeed("RetryAuthor",
		rssEntry("Retry Video", "https://www.youtube.com/watch?v=retry1", daysAgo(1)),
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/atom+xml")
		fmt.Fprint(w, feed)
	}))
	defer srv.Close()

	Update(UpdateOptions{
		URLs:    []string{srv.URL},
		Retries: 2,
		Timeout: 5,
		Fetcher: makeFetcher(srv),
	})

	// The fetcher used in makeFetcher always hits srv.URL directly, so the
	// retry logic in fetchFeedWithRetry will call fetcher again which hits
	// the server a second time and gets a 200. Verify at least one entry inserted.
	if database.TotalEntries() < 1 {
		t.Error("expected at least 1 entry after retry succeeds")
	}
}
