// ---------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Evaluation of Emby DTO & mapping fields to display structures
// ---------------------------------------------------------------------------------------------------------------------

package api

import (
	"EmbyExplorer_for_Windows/models"
	"path/filepath"
	"sort"
	"strconv"
)

const (
	maxActors    = 5
	maxDirectors = 2
	maxStudios   = 1
)

const placeHolder = "-"

func GetFields(collectiontype string) string {
	var m = ""
	switch collectiontype {
	case CollectionMovies:
		m = models.MoviesTableDescription.APIFields
		break
	case CollectionTVShows:
		m = models.SeriesTableDescription.APIFields
		break
	case CollectionHomeVideos:
		m = models.VideosTableDescription.APIFields
		break
	default:
	}
	return m
}

func GetMovieDisplayData(dto []BaseItemDto) []*models.MovieData {
	var result []*models.MovieData
	var movie *models.MovieData
	for _, d := range dto {
		movie = new(models.MovieData)
		movie.MovieId = d.Id
		movie.Name = d.Name
		movie.OriginalTitle = d.OriginalTitle
		movie.ProductionYear = strconv.Itoa(int(d.ProductionYear))
		movie.Studios = evalStudios(d.Studios)
		movie.Actors, movie.Directors = evalPeople(d.People)
		movie.Genres = evalGenres(d.Genres)
		movie.Container = d.Container
		movie.Resolution = evalResolution(d.Width, d.Height)
		movie.Codecs = evalCodecs(d.MediaSources)
		movie.Runtime = evalRuntime(d.RunTimeTicks)
		movie.Path = filepath.Base(d.Path)
		movie.Overview = d.Overview
		result = append(result, movie)
	}
	return result
}

func GetSeriesDisplayData(dto []BaseItemDto) []*models.SeriesData {
	var result []*models.SeriesData
	series := make([]models.SeriesData, 0)
	seasons := make([]models.SeriesData, 0)
	episodes := make([]models.SeriesData, 0)
	var item models.SeriesData
	for _, d := range dto {
		item = models.SeriesData{}
		switch d.Type_ {
		case SeriesType:
			item.Name = d.Name
			item.Actors, _ = evalPeople(d.People)
			item.Genres = evalGenres(d.Genres)
			item.Studios = evalStudios(d.Studios)
			item.Path = d.Path
			item.SeriesId = d.Id
			item.Type_ = d.Type_
			series = append(series, item)
			break
		case SeasonType:
			item.Season = d.Name
			item.SeriesId = d.SeriesId
			item.SeasonId = d.Id
			item.SortIndex = d.IndexNumber
			item.Path = d.Path
			item.Type_ = d.Type_
			seasons = append(seasons, item)
			break
		case EpisodeType:
			item.Episode = d.Name
			item.EpisodeId = d.Id
			item.Runtime = evalRuntime(d.RunTimeTicks)
			item.Container = d.Container
			item.Codecs = evalCodecs(d.MediaSources)
			item.Resolution = evalResolution(d.Width, d.Height)
			item.ProductionYear = strconv.Itoa(int(d.ProductionYear))
			item.Actors, _ = evalPeople(d.People)
			item.SortIndex = d.IndexNumber
			item.Path = filepath.Base(d.Path)
			item.Overview = d.Overview
			item.SeriesId = d.SeriesId
			item.SeasonId = d.SeasonId
			item.Type_ = d.Type_
			episodes = append(episodes, item)
			break
		default:
		}
	}
	// Sort series by Name
	sort.Slice(series, func(i, j int) bool {
		return series[i].Name < series[j].Name
	})
	// Sort seasons by series
	sort.Slice(seasons, func(i, j int) bool {
		return seasons[i].SeriesId < seasons[j].SeriesId
	})
	// Sort episodes by series
	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].SeriesId < episodes[j].SeriesId
	})
	for _, s := range series {
		result = append(result, &s)
		seasonstmp := make([]models.SeriesData, 0)
		// Find seasons for series
		for _, season := range seasons {
			if season.SeriesId == s.SeriesId {
				seasonstmp = append(seasonstmp, season)
			}
		}
		// Sort seasons by IndexNumber
		sort.Slice(seasonstmp, func(i, j int) bool {
			return seasonstmp[i].SortIndex < seasonstmp[j].SortIndex
		})
		for _, n := range seasonstmp {
			// Find episodes for series and season
			episodesstmp := make([]models.SeriesData, 0)
			for _, episode := range episodes {
				if episode.SeriesId == n.SeriesId && episode.SeasonId == n.SeasonId {
					episodesstmp = append(episodesstmp, episode)
				}
			}
			// Sort episodes by IndexNumber
			sort.Slice(episodesstmp, func(i, j int) bool {
				return episodesstmp[i].SortIndex < episodesstmp[j].SortIndex
			})
			for _, e := range episodesstmp {
				e.Name = s.Name
				e.Season = n.Season
				result = append(result, &e)
			}
		}
	}
	return result
}

func GetVideoDisplayData(dto []BaseItemDto) []*models.VideoData {
	var result []*models.VideoData
	folders := make([]models.VideoData, 0)
	videos := make([]models.VideoData, 0)
	var video, folder models.VideoData
	for _, d := range dto {
		switch d.Type_ {
		case VideoType:
			video = models.VideoData{}
			video.Name = d.Name
			video.Container = d.Container
			video.Resolution = evalResolution(d.Width, d.Height)
			video.Codecs = evalCodecs(d.MediaSources)
			video.Runtime = evalRuntime(d.RunTimeTicks)
			video.Path = filepath.Base(d.Path)
			video.ParentId = d.ParentId
			videos = append(videos, video)
			break
		case FolderType:
			folder = models.VideoData{}
			folder.Name = d.Name
			folder.FolderId = d.Id
			folders = append(folders, folder)
			break
		default:
		}
	}
	// Sort folders by Name
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].Name < folders[j].Name
	})
	// Sort videos by Name
	sort.Slice(videos, func(i, j int) bool {
		return videos[i].Name < videos[j].Name
	})
	for _, f := range folders {
		for _, v := range videos {
			if v.ParentId == f.FolderId {
				v.Folder = f.Name
				result = append(result, &v)
			}
		}
	}
	return result
}

func evalStudios(studios []NameLongIdPair) string {
	var s = ""
	for i, studio := range studios {
		if i >= maxStudios {
			break
		}
		s = commaString(s, studio.Name)
	}
	return s
}

func evalPeople(people []BaseItemPerson) (string, string) {
	var actors, directors = "", ""
	var countactors, countdirectors = 0, 0
	for _, p := range people {
		if *p.Type_ == ACTOR_PersonType {
			countactors++
			if countactors <= maxActors {
				actors = commaString(actors, p.Name)
			}
		}
		if *p.Type_ == DIRECTOR_PersonType {
			countdirectors++
			if countdirectors <= maxDirectors {
				directors = commaString(directors, p.Name)
			}
		}
		if countactors > maxActors && countdirectors > maxDirectors {
			break
		}
	}
	return actors, directors
}

func evalGenres(genres []string) string {
	var s = ""
	for _, genre := range genres {
		s = commaString(s, genre)
	}
	return s
}

func evalRuntime(ticks int64) string {
	var s = ""
	if ticks > 0 {
		r := ticks / 10000000
		hours := r / 3600
		minutes := (r % 3600) / 60
		if hours > 0 {
			s = strconv.Itoa(int(hours)) + "h"
		}
		if minutes > 0 {
			s = s + strconv.Itoa(int(minutes)) + "m"
		}
	}
	return s
}

func evalCodecs(media []MediaSourceInfo) string {
	var codecs, codecVideo, codecAudio = "", "", ""
	for _, m := range media {
		for _, s := range m.MediaStreams {
			if *s.Type_ == VIDEO_MediaStreamType {
				codecVideo = s.Codec
			}
			if *s.Type_ == AUDIO_MediaStreamType {
				codecAudio = s.Codec
			}
		}
		if codecVideo == "" {
			codecVideo = placeHolder
		}
		if codecAudio == "" {
			codecAudio = placeHolder
		}
		codecs = codecVideo + ", " + codecAudio
		break
	}
	return codecs
}

func evalResolution(w int32, h int32) string {
	var r = ""
	if w > 0 && h > 0 {
		r = strconv.Itoa(int(w)) + "x" + strconv.Itoa(int(h))
	}
	return r
}

func commaString(source string, append string) string {
	s := source
	if s != "" {
		s = s + ", " + append
	} else {
		s = append
	}
	return s
}
