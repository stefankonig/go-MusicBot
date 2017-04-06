package player

import "testing"

func TestURLParsing(t *testing.T) {
	p := &SpotifyPlayer{}
	tp, id, uid, err := p.getTypeAndIDFromURL("https://open.spotify.com/track/4uLU6hMCjMI75M1A2tKUQC")
	if tp != TYPE_TRACK || id != "4uLU6hMCjMI75M1A2tKUQC" || err != nil {
		t.Log(string(tp), id, uid, err)
		t.Fail()
	}

	tp, id, uid, err = p.getTypeAndIDFromURL("https://open.spotify.com/album/4uLU6hMCjMI75M1A2tKUQC")
	if tp != TYPE_ALBUM || id != "4uLU6hMCjMI75M1A2tKUQC" || err != nil {
		t.Log(string(tp), id, uid, err)
		t.Fail()
	}

	tp, id, uid, err = p.getTypeAndIDFromURL("https://open.spotify.com/user/tana.cross/playlist/2xLFotd9GVVQ6Jde7B3i3B")
	if tp != TYPE_PLAYLIST || id != "2xLFotd9GVVQ6Jde7B3i3B" || err != nil || uid != "tana.cross" {
		t.Log(string(tp), id, uid, err)
		t.Fail()
	}

	tp, id, uid, err = p.getTypeAndIDFromURL("spotify:track:2cBGl1Ehr1D9xbqNmraqb4")
	if tp != TYPE_TRACK || id != "2cBGl1Ehr1D9xbqNmraqb4" || err != nil {
		t.Log(string(tp), id, uid, err)
		t.Fail()
	}

	tp, id, uid, err = p.getTypeAndIDFromURL("spotify:user:111208973:playlist:4XGuyS11n99eMqe1OvN8jq")
	if tp != TYPE_PLAYLIST || id != "4XGuyS11n99eMqe1OvN8jq" || err != nil || uid != "111208973" {
		t.Log(string(tp), id, uid, err)
		t.Fail()
	}
}