package regia

import "github.com/eatmoreapple/regia/internal"

func SetJsonSerializer(serializer internal.Serializer) {
	internal.JSON = serializer
}

func SetXmlSerializer(serializer internal.Serializer) {
	internal.XML = serializer
}
