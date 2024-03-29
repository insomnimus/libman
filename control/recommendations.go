package control

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/zmb3/spotify"
)

type RecommendationEngine uint8

//go:generate stringer -type=RecommendationEngine -trimprefix=Recommendation

const (
	RecommendationNormal RecommendationEngine = iota
	RecommendationExtreme
	RecommendationMin
	RecommendationMax
)

func (r RecommendationEngine) attr(f Features) *spotify.TrackAttributes {
	switch r {
	case RecommendationNormal:
		return f.attrNormal()
	case RecommendationExtreme:
		return f.attrExtreme()
	case RecommendationMin:
		return f.attrMin()
	case RecommendationMax:
		return f.attrMax()
	default:
		panic("internal error: unhandled recommendation engine case")
	}
}

type Features spotify.AudioFeatures

func (a *Features) add(b Features) {
	a.Duration += b.Duration
	a.Acousticness += b.Acousticness
	a.Danceability += b.Danceability
	a.Energy += b.Energy
	a.Instrumentalness += b.Instrumentalness
	a.Liveness += b.Liveness
	a.Loudness += b.Loudness
	a.Speechiness += b.Speechiness
	a.Valence += b.Valence
}

func (a *Features) divide(n float32) {
	a.Duration /= int(n)
	a.Acousticness /= n
	a.Danceability /= n
	a.Energy /= n
	a.Instrumentalness /= n
	a.Liveness /= n
	a.Loudness /= n
	a.Speechiness /= n
	a.Valence /= n
}

// attr converts Features to TrackAttributes.
func (f Features) attrNormal() *spotify.TrackAttributes {
	x := spotify.NewTrackAttributes()
	x.TargetAcousticness(float64(f.Acousticness))
	x.TargetDanceability(float64(f.Danceability))
	x.TargetDuration(f.Duration)
	x.TargetEnergy(float64(f.Energy))
	x.TargetInstrumentalness(float64(f.Instrumentalness))
	x.TargetLiveness(float64(f.Liveness))
	x.TargetLoudness(float64(f.Loudness))
	x.TargetSpeechiness(float64(f.Speechiness))
	x.TargetValence(float64(f.Valence))
	return x
}

// attrExtremes returns recommendations based on fields whose values are more extreme (<25% or >75%).
func (f Features) attrExtreme() *spotify.TrackAttributes {
	ok := func(val float32) bool {
		return val <= 0.25 || val >= 0.75
	}
	x := spotify.NewTrackAttributes()
	if ok(f.Acousticness) {
		x.TargetAcousticness(float64(f.Acousticness))
	}

	if ok(f.Danceability) {
		x.TargetDanceability(float64(f.Danceability))
	}
	x.TargetDuration(f.Duration)
	if ok(f.Energy) {
		x.TargetEnergy(float64(f.Energy))
	}
	if ok(f.Instrumentalness) {
		x.TargetInstrumentalness(float64(f.Instrumentalness))
	}
	if ok(f.Liveness) {
		x.TargetLiveness(float64(f.Liveness))
	}
	if ok(f.Loudness) {
		x.TargetLoudness(float64(f.Loudness))
	}
	if ok(f.Speechiness) {
		x.TargetSpeechiness(float64(f.Speechiness))
	}
	if ok(f.Valence) {
		x.TargetValence(float64(f.Valence))
	}
	return x
}

func (f Features) attrMin() *spotify.TrackAttributes {
	ok := func(val float32) bool {
		return val >= 0.4
	}
	x := spotify.NewTrackAttributes()
	if ok(f.Acousticness) {
		x.MinAcousticness(float64(f.Acousticness))
	}
	if ok(f.Danceability) {
		x.MinDanceability(float64(f.Danceability))
	}
	if ok(f.Energy) {
		x.MinEnergy(float64(f.Energy))
	}
	if ok(f.Instrumentalness) {
		x.MinInstrumentalness(float64(f.Instrumentalness))
	}
	if ok(f.Liveness) {
		x.MinLiveness(float64(f.Liveness))
	}
	if ok(f.Loudness) {
		x.MinLoudness(float64(f.Loudness))
	}
	if ok(f.Speechiness) {
		x.MinSpeechiness(float64(f.Speechiness))
	}
	if ok(f.Valence) {
		x.MinValence(float64(f.Valence))
	}
	return x
}

func (f Features) attrMax() *spotify.TrackAttributes {
	ok := func(val float32) bool {
		return val <= 0.4
	}
	x := spotify.NewTrackAttributes()
	if ok(f.Acousticness) {
		x.MaxAcousticness(float64(f.Acousticness))
	}
	if ok(f.Danceability) {
		x.MaxDanceability(float64(f.Danceability))
	}
	if ok(f.Energy) {
		x.MaxEnergy(float64(f.Energy))
	}
	if ok(f.Instrumentalness) {
		x.MaxInstrumentalness(float64(f.Instrumentalness))
	}
	if ok(f.Liveness) {
		x.MaxLiveness(float64(f.Liveness))
	}
	if ok(f.Loudness) {
		x.MaxLoudness(float64(f.Loudness))
	}
	if ok(f.Speechiness) {
		x.MaxSpeechiness(float64(f.Speechiness))
	}
	if ok(f.Valence) {
		x.MaxValence(float64(f.Valence))
	}
	return x
}

func (p *Playlist) seeds(engine RecommendationEngine) (*spotify.Seeds, *spotify.TrackAttributes, error) {
	if err := p.makeFull(); err != nil {
		return nil, nil, err
	}
	if len(p.Tracks.Tracks) == 0 {
		return nil, nil, fmt.Errorf("%s has no tracks", p.Name)
	}

	var seeds spotify.Seeds
	ids := make([]spotify.ID, len(p.Tracks.Tracks))
	for n, i := range rand.Perm(len(p.Tracks.Tracks)) {
		if n == 5 {
			break
		}
		if rand.Intn(7) == 0 {
			seeds.Artists = append(seeds.Artists, p.Tracks.Tracks[i].Track.Artists[0].ID)
		} else {
			seeds.Tracks = append(seeds.Tracks, p.Tracks.Tracks[i].Track.ID)
		}
	}

	for i, t := range p.Tracks.Tracks {
		ids[i] = t.Track.ID
	}

	feats, err := client.GetAudioFeatures(ids...)
	if err != nil {
		return nil, nil, err
	}

	var n int
	var feat Features
	for _, f := range feats {
		if f != nil {
			n++
			feat.add(Features(*f))
		}
	}

	if n == 0 {
		return &seeds, nil, nil
	}

	feat.divide(float32(n))
	return &seeds, engine.attr(feat), nil
}

func handleRecommend(arg string) error {
	p, engine := chooseRecommendPlaylist(arg)
	if p == nil {
		return nil
	}

	seeds, attr, err := p.seeds(engine)
	if err != nil {
		return err
	}
	recs, err := client.GetRecommendations(*seeds, attr, nil)
	if err != nil {
		return err
	}

	if len(recs.Tracks) == 0 {
		fmt.Println("Spotify returned no recommendations!")
		return nil
	}

	fmt.Println("Playing these recommendations:")
	uris := make([]spotify.URI, len(recs.Tracks))
	for i, t := range recs.Tracks {
		fmt.Printf("%s by %s\n", t.Name, joinArtists(t.Artists))
		uris[i] = t.URI
	}

	err = client.PlayOpt(&spotify.PlayOptions{
		URIs: uris,
	})

	if err == nil {
		isPlaying = true
	}
	return err
}

func chooseRecommendPlaylist(arg string) (*Playlist, RecommendationEngine) {
	if arg == "" {
		eng, cancelled := chooseRecommendationEngine()
		if cancelled {
			return nil, 0
		}
		p := choosePlaylist("")
		return p, eng
	}

	if strings.Contains(arg, "::") {
		split := strings.SplitN(arg, "::", 2)
		left := strings.TrimSpace(split[0])
		var right string
		if len(split) > 1 {
			right = strings.TrimSpace(split[1])
		}
		eng := RecommendationNormal
		switch strings.ToLower(right) {
		case "normal", "":
		case "x", "extreme":
			eng = RecommendationExtreme
		case "min", "minimum":
			eng = RecommendationMin
		case "max", "maximum":
			eng = RecommendationMax
		default:
			fmt.Printf("%s is not a valid recommendation engine.\nRun `help recommend` for the usage.\n", right)
			return nil, 0
		}
		return choosePlaylist(left), eng
	}

	return choosePlaylist(arg), RecommendationNormal
}

func chooseRecommendationEngine() (RecommendationEngine, bool) {
	fmt.Println(`Available recommendation engines:
-	normal
-	extreme
-	min
-	max`)
	for {
		reply, cancelled := readPrompt(false, "Recommendation engine: ")
		if cancelled {
			return 0, true
		}
		switch strings.ToLower(reply) {
		case "normal":
			return RecommendationNormal, false
		case "x", "extreme":
			return RecommendationExtreme, false
		case "min", "minimum":
			return RecommendationMin, false
		case "max", "maximum":
			return RecommendationMax, false
		case "":
			fmt.Println("cancelled")
			return 0, true
		default:
			fmt.Printf("%s is not a valid recommendation engine.\n", reply)
		}
	}
}
