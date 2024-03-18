# `sinister`

`sinister` is a tool to sync and download videos from Youtube.

## Config

```toml
videoFolder = "~/sources"
urls = [
	"https://www.youtube.com/feeds/videos.xml?channel_id=UCeeFfhMcJa1kjtfZAGskOCA",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCdBK94H6oZT2Q7l0-b0xmMg",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UC0vBXGSyV14uvJ4hECDOl0Q",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCHDxYLv8iovIbhrfl16CNyg",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCld68syR8Wi-GY_n4CaoJGA",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCBJycsmduvYEL83R_U4JriQ",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UC0rE2qq81of4fojo-KhO5rg",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCxzC4EngIsMrPmbm6Nxvb-A",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UC6uKrU_WqJ1R2HMTY3LIx5Q",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCBNHHEoiSF8pcLgqLKVugOw",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCHnyfMqiRRG1u-2MsSQLbXA",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCrqM0Ym_NbK1fqeQG2VIohg",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCd3dNckv1Za2coSaHGHl5aA",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCUyeluBRhGPCW4rPe_UvBZQ",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCy0tKL1T7wFoYcxCe0xjN6Q",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCSju5G2aFaWMqn-_0YBtq5A",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCR6LasBpceuYUhuLToKBzvQ",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCb_MAhL8Thb3HJ_wPkH3gcw",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCaSCt8s_4nfkRglWCvNSDrg",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCJfJWct8jN1RpCuVWk3zHTA",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCsBjURrPoezykLs9EqgamOA",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCv1Kcz-CuGM6mxzL3B1_Eiw",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCqJ-Xo29CKyLTjn6z2XwYAw",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UC6me-RzbQFQ-kRyr6BlGZWg",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCbxQcz9k0NRRuy0ukgQTDQQ"
]
```

`sinister` expects RSS feeds of Youtube Channels in it's config file.

You can run `sinister update` to update the database.

Then you can run `sinister download`, to download a unwatched video.



