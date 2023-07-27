package runner

import (
	"encoding/csv"
	"fmt"
	"github.com/iahfdoa/m3u8Find/pkg/result"
	"github.com/iahfdoa/m3u8Find/pkg/utils"
	"github.com/pkg/errors"
	"github.com/projectdiscovery/gologger"
	"io"
	"os"
	"reflect"
)

// WriteCsvOutput writes the output list of subdomain in csv format to an io.Writer
func writeCsvOutput(data *result.UsersInfo, header bool, writer io.Writer) error {
	encoder := csv.NewWriter(writer)
	if header {
		writeCSVHeaders(data, encoder)
	}
	writeCSVRow(data, encoder)

	encoder.Flush()
	return nil
}
func writeCSVHeaders(data *result.UsersInfo, writer *csv.Writer) {
	headers, err := data.CSVHeaders()
	if err != nil {
		gologger.Error().Msgf(err.Error())
		return
	}

	if err := writer.Write(headers); err != nil {
		errMsg := errors.Wrap(err, "Could not write headers")
		gologger.Error().Msgf(errMsg.Error())
	}
}
func writeCSVRow(data *result.UsersInfo, writer *csv.Writer) {
	rowData, err := data.CSVFields()
	if err != nil {
		gologger.Error().Msgf(err.Error())
		return
	}
	if err := writer.Write(rowData); err != nil {
		errMsg := errors.Wrap(err, "Could not write row")
		gologger.Error().Msgf(errMsg.Error())
	}
}
func setM3u8Header(s string) string {
	s = "#EXTM3U x-tvg-url = \"\"\n" + s
	return s
}
func outputCsvAllUsers(filename string, us []result.UsersInfo) error {
	if us == nil {
		return errors.New("不存在主播")
	}
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		if os.IsPermission(err) {
			return errors.Wrap(err, "permission error")
		}
		return errors.New("打开文件失败: " + err.Error())
	}

	for i, u := range us {
		header := false
		if i == 0 {
			header = true
		}
		err := writeCsvOutput(&u, header, f)
		if err != nil {
			return err
		}
	}
	return nil
}
func outputM3u8Users(filename, groupRole string, users []result.UsersInfo) error {
	if users == nil {
		return errors.New("沒有主播信息")
	}
	var contents string
	contents = setM3u8Header("")
	for _, u := range users {
		tvgId := u.Id
		tvgName := u.Username
		grT := reflect.ValueOf(u).FieldByName(utils.FirstUpper(groupRole))
		groupTitle := grT.String()
		tvgLog := u.PreviewUrlThumbSmall
		title := u.Username
		m3u8 := u.HlsPlaylist
		contents += fmt.Sprintf("# EXTINF:-1, tvg-id = \"%v\",tvg-name=\"%v\",group-title = \"%v\", tvg-logo = \"%v\",%v\n%v\n", tvgId, tvgName, groupTitle, tvgLog, title, m3u8)
	}
	err := utils.WriteFile(contents, filename)
	if err != nil {
		return err
	}
	return nil
}
func outputM3u8AllUsers(filename string, users *[]result.UsersInfo) error {
	if users == nil {
		return errors.New("沒有主播信息")
	}
	var contents string
	contents = setM3u8Header("")
	for _, u := range *users {
		tvgId := u.Id
		tvgName := u.Username
		tvgLog := u.PreviewUrlThumbSmall
		groupTitle := u.ParentTag
		title := u.Username
		m3u8 := u.HlsPlaylist
		contents += fmt.Sprintf("# EXTINF:-1, tvg-id = \"%v\",tvg-name=\"%v\",group-title = \"%v\", tvg-logo = \"%v\",%v\n%v\n", tvgId, tvgName, groupTitle, tvgLog, title, m3u8)
	}

	err := utils.WriteFile(contents, filename)
	if err != nil {
		return err
	}
	return nil
}
