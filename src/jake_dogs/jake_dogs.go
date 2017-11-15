package jake_dogs

import (
    "math/rand"
)

var Images []string

func Init() {
    Images = []string{
        "https://pbs.twimg.com/media/CAhcKu6U0AA64py.jpg:large",
        "http://78.media.tumblr.com/1f35ceb6c7bdb71d243e1a12d48bcc56/tumblr_nlfp27h6sw1rv5rhio3_250.jpg",
        "http://78.media.tumblr.com/0c6256ebcbd3003b5ccecc24fb3242ca/tumblr_nlfp27h6sw1rv5rhio1_400.jpg",
        "http://78.media.tumblr.com/65df90f4c36ba0579fcf7a01c60b8b65/tumblr_nlfp27h6sw1rv5rhio4_1280.jpg",
    }
}

func GetImage() string {
    ret_url := Images[rand.Int() % len(Images)]
    return ret_url
}