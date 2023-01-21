package utils

import (
	"testing"

	"github.com/spf13/viper"
)

func TestCheckCharLimit(t *testing.T) {
	viper.AutomaticEnv()
	viper.AddConfigPath("../")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.Set("PRIVATEBIN_ENABLE", false)
	msg := `
	Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. 
	Dictum varius duis at consectetur lorem donec massa sapien faucibus. Rhoncus est pellentesque elit ullamcorper dignissim. 
	Elit sed vulputate mi sit amet mauris commodo quis. Sed euismod nisi porta lorem mollis. 
	Id velit ut tortor pretium viverra suspendisse potenti nullam ac. Felis imperdiet proin fermentum leo vel orci. In iaculis nunc sed augue lacus viverra vitae congue eu. 
	Ornare quam viverra orci sagittis eu volutpat odio facilisis mauris. Ut consequat semper viverra nam libero justo laoreet. 
	Ac tortor dignissim convallis aenean et. Integer vitae justo eget magna fermentum. Urna cursus eget nunc scelerisque viverra mauris. 
	Luctus venenatis lectus magna fringilla urna porttitor rhoncus dolor.
	Nullam eget felis eget nunc lobortis. Fermentum posuere urna nec tincidunt praesent semper. Quis eleifend quam adipiscing vitae proin sagittis nisl. 
	Scelerisque mauris pellentesque pulvinar pellentesque. Lorem sed risus ultricies tristique nulla. Eget arcu dictum varius duis at consectetur lorem. 
	Eleifend mi in nulla posuere sollicitudin aliquam ultrices sagittis.
	`
	limit := 155
	msg, err := CheckCharLimit(msg, limit, false)
	got := len(msg)
	want := 188 // 155+33 (\nMessage limit exceeded.``````\n)
	if err != nil {
		t.Errorf(msg)
	}
	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}
