package main

type Transport struct{}

type Config struct {
	SonySelectID  string
	OutPath       string
	TrackTemplate string
	OmitArtists   bool
	KeepCover     bool
	Urls          []string
}

type Args struct {
	Urls    []string `arg:"positional, required"`
	OutPath string   `arg:"-o" help:"Where to download to. Path will be made if it doesn't already exist."`
}

type PostData struct {
	Content struct {
		AlbumID string `json:"albumId"`
		MusicId int    `json:"musicId"`
		IndexID int    `json:"indexId"`
	} `json:"content"`
	Header struct {
		AccessKey         string `json:"accessKey"`
		ContentEncryption bool   `json:"contentEncryption"`
		Imei              string `json:"imei"`
		Model             string `json:"model"`
		Nonce             string `json:"nonce"`
		SignEncryption    bool   `json:"signEncryption"`
		SonySelectID      string `json:"sonySelectId"`
		Timestamp         int64  `json:"timestamp"`
		Version           string `json:"version"`
	} `json:"header"`
}

type AlbumMetaPost struct {
	Content struct {
		AlbumID string `json:"albumId"`
	} `json:"content"`
	Header struct {
		AccessKey         string `json:"accessKey"`
		ContentEncryption bool   `json:"contentEncryption"`
		Imei              string `json:"imei"`
		Model             string `json:"model"`
		Nonce             string `json:"nonce"`
		SignEncryption    bool   `json:"signEncryption"`
		SonySelectID      string `json:"sonySelectId"`
		Timestamp         int64  `json:"timestamp"`
		Version           string `json:"version"`
	} `json:"header"`
}

type TrackMetaPost struct {
	Content struct {
		MusicId int `json:"musicId"`
	} `json:"content"`
	Header struct {
		AccessKey         string `json:"accessKey"`
		ContentEncryption bool   `json:"contentEncryption"`
		Imei              string `json:"imei"`
		Model             string `json:"model"`
		Nonce             string `json:"nonce"`
		SignEncryption    bool   `json:"signEncryption"`
		SonySelectID      string `json:"sonySelectId"`
		Timestamp         int64  `json:"timestamp"`
		Version           string `json:"version"`
	} `json:"header"`
}

type FileMetaPost struct {
	Content struct {
		IndexID int `json:"indexId"`
	} `json:"content"`
	Header struct {
		AccessKey         string `json:"accessKey"`
		ContentEncryption bool   `json:"contentEncryption"`
		Imei              string `json:"imei"`
		Model             string `json:"model"`
		Nonce             string `json:"nonce"`
		SignEncryption    bool   `json:"signEncryption"`
		SonySelectID      string `json:"sonySelectId"`
		Timestamp         int64  `json:"timestamp"`
		Version           string `json:"version"`
	} `json:"header"`
}

type AlbumMeta struct {
	Result struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
	} `json:"result"`
	Content struct {
		ReleaseTime   string  `json:"releaseTime"`
		Singer        string  `json:"singer"`
		Artist        string  `json:"artist"`
		IsHiRes       bool    `json:"isHiRes"`
		DiscountPrice float64 `json:"discountPrice"`
		AlbumID       int     `json:"albumId"`
		Description   string  `json:"description"`
		Bitrate       string  `json:"bitrate"`
		CdList        []struct {
			Composer  string `json:"composer"`
			Type      string `json:"type"`
			WorkName  string `json:"workName"`
			Musiclist []struct {
				Duration  string  `json:"duration"`
				MusicID   int     `json:"musicId"`
				Size      float64 `json:"size"`
				Artist    string  `json:"artist"`
				Composer  string  `json:"composer"`
				Property  string  `json:"property"`
				TrackNo   int     `json:"trackNo"`
				MusicName string  `json:"musicName"`
				Isfollow  bool    `json:"isfollow"`
				Promotion string  `json:"promotion"`
			} `json:"musiclist"`
		} `json:"cdList"`
		ListCa []struct {
			CategoryType string `json:"categoryType"`
			CategoryName string `json:"categoryName"`
			CategoryID   int    `json:"categoryId"`
		} `json:"listCa"`
		Price         float64     `json:"price"`
		BackCover     interface{} `json:"backCover"`
		SmallIcon     string      `json:"smallIcon"`
		Brand         string      `json:"brand"`
		Player        string      `json:"player"`
		CommentNumber string      `json:"commentNumber"`
		Resource      struct {
			Image     []interface{} `json:"image"`
			Pdf       []interface{} `json:"pdf"`
			Video     []interface{} `json:"video"`
			Interview []interface{} `json:"interview"`
		} `json:"resource"`
		Composer   string `json:"composer"`
		Format     string `json:"format"`
		PlayModels []struct {
			Size       float64 `json:"size"`
			Property   string  `json:"property,omitempty"`
			Format     string  `json:"format"`
			AlbumID    int     `json:"albumId"`
			Bitrate    string  `json:"bitrate"`
			Permission bool    `json:"permission"`
			Type       string  `json:"type"`
			Selected   bool    `json:"selected"`
		} `json:"playModels"`
		HasComment int     `json:"hasComment"`
		Conductor  string  `json:"conductor"`
		IsFollow   int     `json:"isFollow"`
		Size       float64 `json:"size"`
		Name       string  `json:"name"`
		LargeIcon  string  `json:"largeIcon"`
	} `json:"content"`
	ContentEncryption bool `json:"contentEncryption"`
}

type TrackMeta struct {
	Result struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
	} `json:"result"`
	Content struct {
		AlbumName   string `json:"albumName"`
		ReleaseTime string `json:"releaseTime"`
		Singer      string `json:"singer"`
		Artist      string `json:"artist"`
		IsHiRes     bool   `json:"isHiRes"`
		Icon        string `json:"icon"`
		AlbumID     int    `json:"albumId"`
		Description string `json:"description"`
		Bitrate     string `json:"bitrate"`
		Duration    int    `json:"duration"`
		MusicID     int    `json:"musicId"`
		ListCa      []struct {
			CategoryType string `json:"categoryType"`
			CategoryName string `json:"categoryName"`
			CategoryID   int    `json:"categoryId"`
		} `json:"listCa"`
		Brand         string `json:"brand"`
		Player        string `json:"player"`
		CommentNumber string `json:"commentNumber"`
		Is360RA       bool   `json:"is360RA"`
		Composer      string `json:"composer"`
		HasFollow     int    `json:"hasFollow"`
		PlayModels    []struct {
			MusicID    int     `json:"musicId"`
			Size       float64 `json:"size"`
			Format     string  `json:"format"`
			AlbumID    int     `json:"albumId"`
			IndexID    int     `json:"indexId"`
			Bitrate    string  `json:"bitrate"`
			Permission bool    `json:"permission"`
			Type       string  `json:"type"`
			Selected   bool    `json:"selected"`
		} `json:"playModels"`
		FollowNumber string   `json:"followNumber"`
		WorkName     string   `json:"workName"`
		HasComment   int      `json:"hasComment"`
		Conductor    string   `json:"conductor"`
		ListRate     []string `json:"listRate"`
		MusicName    string   `json:"musicName"`
		IsFollow     int      `json:"isFollow"`
		Corporation  struct {
			Name interface{} `json:"name"`
			ID   int         `json:"id"`
		} `json:"Corporation"`
		Size       float64  `json:"size"`
		ListFormat []string `json:"listFormat"`
	} `json:"content"`
	ContentEncryption bool `json:"contentEncryption"`
}

type FileMetaEnc struct {
	Result struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
	} `json:"result"`
	Content struct {
		Encrypcontent string `json:"encrypcontent"`
	} `json:"content"`
	ContentEncryption bool `json:"contentEncryption"`
}

type FileMeta struct {
	SecretKey           string   `json:"secretKey"`
	SegmentSize         int      `json:"segmentSize"`
	Format              string   `json:"format"`
	AlbumID             int      `json:"albumId"`
	SampleRate          int      `json:"sampleRate"`
	EncryptionAlgorithm string   `json:"encryptionAlgorithm"`
	Samples             int      `json:"samples"`
	Segments            int      `json:"segments"`
	BaseURL             string   `json:"baseUrl"`
	FrameSize           int      `json:"frameSize"`
	Names               []string `json:"names"`
	MusicID             int      `json:"musicId"`
	SampleBit           int      `json:"sampleBit"`
	Property            string   `json:"property"`
}
