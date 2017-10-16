package dash

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"aqwari.net/xml/xmltree"
)

func TestDash(t *testing.T) {
	type Document struct {
		MPD MPDtype `xml:"urn:mpeg:dash:schema:mpd:2011 MPD"`
	}
	var document Document
	samples, err := filepath.Glob(filepath.Join("*.xml"))
	if err != nil {
		t.Fatal(err)
	}
	if len(samples) != 1 {
		t.Fatal("expected one sample file, found ", samples)
	}
	input, err := ioutil.ReadFile(samples[0])
	if err != nil {
		t.Fatal(err)
	}
	input = append([]byte("<Document>\n"), input...)
	input = append(input, []byte("</Document>")...)
	if err := xml.Unmarshal(input, &document); err != nil {
		t.Fatal("unmarshal: ", err)
	}
	output, err := xml.Marshal(&document)
	if err != nil {
		t.Fatal("marshal: ", err)
	}
	inputTree, err := xmltree.Parse(input)
	if err != nil {
		t.Fatal("dash: ", err)
	}
	outputTree, err := xmltree.Parse(output)
	if err != nil {
		t.Fatal("remarshal: ", err)
	}
	if !xmltree.Equal(inputTree, outputTree) {
		t.Errorf("got \n%s\n, wanted \n%s\n", xmltree.MarshalIndent(outputTree, "", "  "), xmltree.MarshalIndent(inputTree, "", "  "))
	}
}

// May be one of onLoad, onRequest, other, none
type ActuateType string

type AdaptationSetType struct {
	RepresentationBaseType
	Actuate                 ActuateType            `xml:"actuate,attr,omitempty"`
	Group                   uint                   `xml:"group,attr,omitempty"`
	Lang                    string                 `xml:"lang,attr,omitempty"`
	ContentType             string                 `xml:"contentType,attr,omitempty"`
	Par                     RatioType              `xml:"par,attr,omitempty"`
	MinBandwidth            uint                   `xml:"minBandwidth,attr,omitempty"`
	MaxBandwidth            uint                   `xml:"maxBandwidth,attr,omitempty"`
	MinWidth                uint                   `xml:"minWidth,attr,omitempty"`
	MaxWidth                uint                   `xml:"maxWidth,attr,omitempty"`
	MinHeight               uint                   `xml:"minHeight,attr,omitempty"`
	MaxHeight               uint                   `xml:"maxHeight,attr,omitempty"`
	MinFrameRate            FrameRateType          `xml:"minFrameRate,attr,omitempty"`
	MaxFrameRate            FrameRateType          `xml:"maxFrameRate,attr,omitempty"`
	SegmentAlignment        ConditionalUintType    `xml:"segmentAlignment,attr,omitempty"`
	SubsegmentAlignment     ConditionalUintType    `xml:"subsegmentAlignment,attr,omitempty"`
	SubsegmentStartsWithSAP uint                   `xml:"subsegmentStartsWithSAP,attr,omitempty"`
	BitstreamSwitching      bool                   `xml:"bitstreamSwitching,attr,omitempty"`
	Accessibility           []DescriptorType       `xml:"urn:mpeg:dash:schema:mpd:2011 Accessibility,omitempty"`
	Role                    []DescriptorType       `xml:"urn:mpeg:dash:schema:mpd:2011 Role,omitempty"`
	Rating                  []DescriptorType       `xml:"urn:mpeg:dash:schema:mpd:2011 Rating,omitempty"`
	Viewpoint               []DescriptorType       `xml:"urn:mpeg:dash:schema:mpd:2011 Viewpoint,omitempty"`
	ContentComponent        []ContentComponentType `xml:"urn:mpeg:dash:schema:mpd:2011 ContentComponent,omitempty"`
	BaseURL                 []BaseURLType          `xml:"urn:mpeg:dash:schema:mpd:2011 BaseURL,omitempty"`
	SegmentBase             SegmentBaseType        `xml:"urn:mpeg:dash:schema:mpd:2011 SegmentBase,omitempty"`
	SegmentList             SegmentListType        `xml:"urn:mpeg:dash:schema:mpd:2011 SegmentList,omitempty"`
	SegmentTemplate         SegmentTemplateType    `xml:"urn:mpeg:dash:schema:mpd:2011 SegmentTemplate,omitempty"`
	Representation          []RepresentationType   `xml:"urn:mpeg:dash:schema:mpd:2011 Representation,omitempty"`
}

func (t *AdaptationSetType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type T AdaptationSetType
	var overlay struct {
		*T
		Actuate                 ActuateType         `xml:"actuate,attr,omitempty"`
		SegmentAlignment        ConditionalUintType `xml:"segmentAlignment,attr,omitempty"`
		SubsegmentAlignment     ConditionalUintType `xml:"subsegmentAlignment,attr,omitempty"`
		SubsegmentStartsWithSAP uint                `xml:"subsegmentStartsWithSAP,attr,omitempty"`
	}
	overlay.T = (*T)(t)
	if err := d.DecodeElement(&overlay, &start); err != nil {
		return err
	}
	overlay.T.Actuate = ActuateType(overlay.Actuate)
	overlay.T.SegmentAlignment = ConditionalUintType(overlay.SegmentAlignment)
	overlay.T.SubsegmentAlignment = ConditionalUintType(overlay.SubsegmentAlignment)
	overlay.T.SubsegmentStartsWithSAP = uint(overlay.SubsegmentStartsWithSAP)
	return nil
}

type BaseURLType struct {
	Value                    string  `xml:",chardata"`
	ServiceLocation          string  `xml:"serviceLocation,attr,omitempty"`
	ByteRange                string  `xml:"byteRange,attr,omitempty"`
	AvailabilityTimeOffset   float64 `xml:"availabilityTimeOffset,attr,omitempty"`
	AvailabilityTimeComplete bool    `xml:"availabilityTimeComplete,attr,omitempty"`
}

type ConditionalUintType string

type ContentComponentType struct {
	Lang          string           `xml:"lang,attr,omitempty"`
	ContentType   string           `xml:"contentType,attr,omitempty"`
	Par           RatioType        `xml:"par,attr,omitempty"`
	Items         []string         `xml:",any"`
	Accessibility []DescriptorType `xml:"urn:mpeg:dash:schema:mpd:2011 Accessibility,omitempty"`
	Role          []DescriptorType `xml:"urn:mpeg:dash:schema:mpd:2011 Role,omitempty"`
	Rating        []DescriptorType `xml:"urn:mpeg:dash:schema:mpd:2011 Rating,omitempty"`
	Viewpoint     []DescriptorType `xml:"urn:mpeg:dash:schema:mpd:2011 Viewpoint,omitempty"`
}

type DescriptorType struct {
	SchemeIdUri string   `xml:"schemeIdUri,attr"`
	Value       string   `xml:"value,attr,omitempty"`
	Items       []string `xml:",any"`
}

type EventStreamType struct {
	Actuate     ActuateType `xml:"actuate,attr,omitempty"`
	SchemeIdUri string      `xml:"schemeIdUri,attr"`
	Value       string      `xml:"value,attr,omitempty"`
	Timescale   uint        `xml:"timescale,attr,omitempty"`
	Items       []string    `xml:",any"`
	Event       []EventType `xml:"urn:mpeg:dash:schema:mpd:2011 Event,omitempty"`
}

func (t *EventStreamType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type T EventStreamType
	var overlay struct {
		*T
		Actuate ActuateType `xml:"actuate,attr,omitempty"`
	}
	overlay.T = (*T)(t)
	if err := d.DecodeElement(&overlay, &start); err != nil {
		return err
	}
	overlay.T.Actuate = ActuateType(overlay.Actuate)
	return nil
}

type EventType struct {
	PresentationTime uint64   `xml:"presentationTime,attr,omitempty"`
	Duration         uint64   `xml:"duration,attr,omitempty"`
	MessageData      string   `xml:"messageData,attr,omitempty"`
	Items            []string `xml:",any"`
}

func (t *EventType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type T EventType
	var overlay struct {
		*T
		PresentationTime uint64 `xml:"presentationTime,attr,omitempty"`
	}
	overlay.T = (*T)(t)
	if err := d.DecodeElement(&overlay, &start); err != nil {
		return err
	}
	overlay.T.PresentationTime = uint64(overlay.PresentationTime)
	return nil
}

// Must match the pattern [0-9]*[0-9](/[0-9]*[0-9])?
type FrameRateType string

type MPDtype struct {
	Profiles                   string                   `xml:"profiles,attr"`
	Type                       PresentationType         `xml:"type,attr,omitempty"`
	AvailabilityStartTime      time.Time                `xml:"availabilityStartTime,attr,omitempty"`
	AvailabilityEndTime        time.Time                `xml:"availabilityEndTime,attr,omitempty"`
	PublishTime                time.Time                `xml:"publishTime,attr,omitempty"`
	MediaPresentationDuration  string                   `xml:"mediaPresentationDuration,attr,omitempty"`
	MinimumUpdatePeriod        string                   `xml:"minimumUpdatePeriod,attr,omitempty"`
	MinBufferTime              string                   `xml:"minBufferTime,attr"`
	TimeShiftBufferDepth       string                   `xml:"timeShiftBufferDepth,attr,omitempty"`
	SuggestedPresentationDelay string                   `xml:"suggestedPresentationDelay,attr,omitempty"`
	MaxSegmentDuration         string                   `xml:"maxSegmentDuration,attr,omitempty"`
	MaxSubsegmentDuration      string                   `xml:"maxSubsegmentDuration,attr,omitempty"`
	Items                      []string                 `xml:",any"`
	ProgramInformation         []ProgramInformationType `xml:"urn:mpeg:dash:schema:mpd:2011 ProgramInformation,omitempty"`
	BaseURL                    []BaseURLType            `xml:"urn:mpeg:dash:schema:mpd:2011 BaseURL,omitempty"`
	Location                   []string                 `xml:"urn:mpeg:dash:schema:mpd:2011 Location,omitempty"`
	Period                     []PeriodType             `xml:"urn:mpeg:dash:schema:mpd:2011 Period"`
	Metrics                    []MetricsType            `xml:"urn:mpeg:dash:schema:mpd:2011 Metrics,omitempty"`
	EssentialProperty          []DescriptorType         `xml:"urn:mpeg:dash:schema:mpd:2011 EssentialProperty,omitempty"`
	SupplementalProperty       []DescriptorType         `xml:"urn:mpeg:dash:schema:mpd:2011 SupplementalProperty,omitempty"`
	UTCTiming                  []DescriptorType         `xml:"urn:mpeg:dash:schema:mpd:2011 UTCTiming,omitempty"`
}

func (t *MPDtype) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type T MPDtype
	var layout struct {
		*T
		AvailabilityStartTime xsdDateTime `xml:"availabilityStartTime,attr,omitempty"`
		AvailabilityEndTime   xsdDateTime `xml:"availabilityEndTime,attr,omitempty"`
		PublishTime           xsdDateTime `xml:"publishTime,attr,omitempty"`
	}
	layout.T = (*T)(t)
	layout.AvailabilityStartTime = xsdDateTime(layout.T.AvailabilityStartTime)
	layout.AvailabilityEndTime = xsdDateTime(layout.T.AvailabilityEndTime)
	layout.PublishTime = xsdDateTime(layout.T.PublishTime)
	return e.EncodeElement(layout, start)
}
func (t *MPDtype) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type T MPDtype
	var overlay struct {
		*T
		Type                  PresentationType `xml:"type,attr,omitempty"`
		AvailabilityStartTime xsdDateTime      `xml:"availabilityStartTime,attr,omitempty"`
		AvailabilityEndTime   xsdDateTime      `xml:"availabilityEndTime,attr,omitempty"`
		PublishTime           xsdDateTime      `xml:"publishTime,attr,omitempty"`
	}
	overlay.T = (*T)(t)
	if err := d.DecodeElement(&overlay, &start); err != nil {
		return err
	}
	overlay.T.Type = PresentationType(overlay.Type)
	overlay.T.AvailabilityStartTime = time.Time(overlay.AvailabilityStartTime)
	overlay.T.AvailabilityEndTime = time.Time(overlay.AvailabilityEndTime)
	overlay.T.PublishTime = time.Time(overlay.PublishTime)
	return nil
}

type MetricsType struct {
	Metrics   string           `xml:"metrics,attr"`
	Items     []string         `xml:",any"`
	Reporting []DescriptorType `xml:"urn:mpeg:dash:schema:mpd:2011 Reporting"`
	Range     []RangeType      `xml:"urn:mpeg:dash:schema:mpd:2011 Range,omitempty"`
}

type MultipleSegmentBaseType struct {
	SegmentBaseType
	Duration           uint                `xml:"duration,attr,omitempty"`
	StartNumber        uint                `xml:"startNumber,attr,omitempty"`
	SegmentTimeline    SegmentTimelineType `xml:"urn:mpeg:dash:schema:mpd:2011 SegmentTimeline,omitempty"`
	BitstreamSwitching URLType             `xml:"urn:mpeg:dash:schema:mpd:2011 BitstreamSwitching,omitempty"`
}

type PeriodType struct {
	Actuate              ActuateType         `xml:"actuate,attr,omitempty"`
	Start                string              `xml:"start,attr,omitempty"`
	Duration             string              `xml:"duration,attr,omitempty"`
	BitstreamSwitching   bool                `xml:"bitstreamSwitching,attr,omitempty"`
	Items                []string            `xml:",any"`
	BaseURL              []BaseURLType       `xml:"urn:mpeg:dash:schema:mpd:2011 BaseURL,omitempty"`
	SegmentBase          SegmentBaseType     `xml:"urn:mpeg:dash:schema:mpd:2011 SegmentBase,omitempty"`
	SegmentList          SegmentListType     `xml:"urn:mpeg:dash:schema:mpd:2011 SegmentList,omitempty"`
	SegmentTemplate      SegmentTemplateType `xml:"urn:mpeg:dash:schema:mpd:2011 SegmentTemplate,omitempty"`
	AssetIdentifier      DescriptorType      `xml:"urn:mpeg:dash:schema:mpd:2011 AssetIdentifier,omitempty"`
	EventStream          []EventStreamType   `xml:"urn:mpeg:dash:schema:mpd:2011 EventStream,omitempty"`
	AdaptationSet        []AdaptationSetType `xml:"urn:mpeg:dash:schema:mpd:2011 AdaptationSet,omitempty"`
	Subset               []SubsetType        `xml:"urn:mpeg:dash:schema:mpd:2011 Subset,omitempty"`
	SupplementalProperty []DescriptorType    `xml:"urn:mpeg:dash:schema:mpd:2011 SupplementalProperty,omitempty"`
}

func (t *PeriodType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type T PeriodType
	var overlay struct {
		*T
		Actuate            ActuateType `xml:"actuate,attr,omitempty"`
		BitstreamSwitching bool        `xml:"bitstreamSwitching,attr,omitempty"`
	}
	overlay.T = (*T)(t)
	if err := d.DecodeElement(&overlay, &start); err != nil {
		return err
	}
	overlay.T.Actuate = ActuateType(overlay.Actuate)
	overlay.T.BitstreamSwitching = bool(overlay.BitstreamSwitching)
	return nil
}

// May be one of static, dynamic
type PresentationType string

type ProgramInformationType struct {
	Lang               string   `xml:"lang,attr,omitempty"`
	MoreInformationURL string   `xml:"moreInformationURL,attr,omitempty"`
	Items              []string `xml:",any"`
	Title              string   `xml:"urn:mpeg:dash:schema:mpd:2011 Title,omitempty"`
	Source             string   `xml:"urn:mpeg:dash:schema:mpd:2011 Source,omitempty"`
	Copyright          string   `xml:"urn:mpeg:dash:schema:mpd:2011 Copyright,omitempty"`
}

type RangeType struct {
	Starttime string `xml:"starttime,attr,omitempty"`
	Duration  string `xml:"duration,attr,omitempty"`
}

// Must match the pattern [0-9]*:[0-9]*
type RatioType string

type RepresentationBaseType struct {
	Profiles                  string            `xml:"profiles,attr,omitempty"`
	Width                     uint              `xml:"width,attr,omitempty"`
	Height                    uint              `xml:"height,attr,omitempty"`
	Sar                       RatioType         `xml:"sar,attr,omitempty"`
	FrameRate                 FrameRateType     `xml:"frameRate,attr,omitempty"`
	AudioSamplingRate         string            `xml:"audioSamplingRate,attr,omitempty"`
	MimeType                  string            `xml:"mimeType,attr,omitempty"`
	SegmentProfiles           string            `xml:"segmentProfiles,attr,omitempty"`
	Codecs                    string            `xml:"codecs,attr,omitempty"`
	MaximumSAPPeriod          float64           `xml:"maximumSAPPeriod,attr,omitempty"`
	StartWithSAP              uint              `xml:"startWithSAP,attr,omitempty"`
	MaxPlayoutRate            float64           `xml:"maxPlayoutRate,attr,omitempty"`
	CodingDependency          bool              `xml:"codingDependency,attr,omitempty"`
	ScanType                  VideoScanType     `xml:"scanType,attr,omitempty"`
	Items                     []string          `xml:",any"`
	FramePacking              []DescriptorType  `xml:"urn:mpeg:dash:schema:mpd:2011 FramePacking,omitempty"`
	AudioChannelConfiguration []DescriptorType  `xml:"urn:mpeg:dash:schema:mpd:2011 AudioChannelConfiguration,omitempty"`
	ContentProtection         []DescriptorType  `xml:"urn:mpeg:dash:schema:mpd:2011 ContentProtection,omitempty"`
	EssentialProperty         []DescriptorType  `xml:"urn:mpeg:dash:schema:mpd:2011 EssentialProperty,omitempty"`
	SupplementalProperty      []DescriptorType  `xml:"urn:mpeg:dash:schema:mpd:2011 SupplementalProperty,omitempty"`
	InbandEventStream         []EventStreamType `xml:"urn:mpeg:dash:schema:mpd:2011 InbandEventStream,omitempty"`
}

type RepresentationType struct {
	RepresentationBaseType
	Bandwidth              uint                    `xml:"bandwidth,attr"`
	QualityRanking         uint                    `xml:"qualityRanking,attr,omitempty"`
	DependencyId           StringVectorType        `xml:"dependencyId,attr,omitempty"`
	MediaStreamStructureId StringVectorType        `xml:"mediaStreamStructureId,attr,omitempty"`
	BaseURL                []BaseURLType           `xml:"urn:mpeg:dash:schema:mpd:2011 BaseURL,omitempty"`
	SubRepresentation      []SubRepresentationType `xml:"urn:mpeg:dash:schema:mpd:2011 SubRepresentation,omitempty"`
	SegmentBase            SegmentBaseType         `xml:"urn:mpeg:dash:schema:mpd:2011 SegmentBase,omitempty"`
	SegmentList            SegmentListType         `xml:"urn:mpeg:dash:schema:mpd:2011 SegmentList,omitempty"`
	SegmentTemplate        SegmentTemplateType     `xml:"urn:mpeg:dash:schema:mpd:2011 SegmentTemplate,omitempty"`
}

type S struct {
	T uint64 `xml:"t,attr,omitempty"`
	N uint64 `xml:"n,attr,omitempty"`
	D uint64 `xml:"d,attr"`
	R int    `xml:"r,attr,omitempty"`
}

func (t *S) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type T S
	var overlay struct {
		*T
		R int `xml:"r,attr,omitempty"`
	}
	overlay.T = (*T)(t)
	if err := d.DecodeElement(&overlay, &start); err != nil {
		return err
	}
	overlay.T.R = int(overlay.R)
	return nil
}

type SegmentBaseType struct {
	Timescale                uint     `xml:"timescale,attr,omitempty"`
	PresentationTimeOffset   uint64   `xml:"presentationTimeOffset,attr,omitempty"`
	IndexRange               string   `xml:"indexRange,attr,omitempty"`
	IndexRangeExact          bool     `xml:"indexRangeExact,attr,omitempty"`
	AvailabilityTimeOffset   float64  `xml:"availabilityTimeOffset,attr,omitempty"`
	AvailabilityTimeComplete bool     `xml:"availabilityTimeComplete,attr,omitempty"`
	Items                    []string `xml:",any"`
	Initialization           URLType  `xml:"urn:mpeg:dash:schema:mpd:2011 Initialization,omitempty"`
	RepresentationIndex      URLType  `xml:"urn:mpeg:dash:schema:mpd:2011 RepresentationIndex,omitempty"`
}

func (t *SegmentBaseType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type T SegmentBaseType
	var overlay struct {
		*T
		IndexRangeExact bool `xml:"indexRangeExact,attr,omitempty"`
	}
	overlay.T = (*T)(t)
	if err := d.DecodeElement(&overlay, &start); err != nil {
		return err
	}
	overlay.T.IndexRangeExact = bool(overlay.IndexRangeExact)
	return nil
}

type SegmentListType struct {
	MultipleSegmentBaseType
	Actuate    ActuateType      `xml:"actuate,attr,omitempty"`
	SegmentURL []SegmentURLType `xml:"urn:mpeg:dash:schema:mpd:2011 SegmentURL,omitempty"`
}

func (t *SegmentListType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type T SegmentListType
	var overlay struct {
		*T
		Actuate ActuateType `xml:"actuate,attr,omitempty"`
	}
	overlay.T = (*T)(t)
	if err := d.DecodeElement(&overlay, &start); err != nil {
		return err
	}
	overlay.T.Actuate = ActuateType(overlay.Actuate)
	return nil
}

type SegmentTemplateType struct {
	MultipleSegmentBaseType
	Media              string `xml:"media,attr,omitempty"`
	Index              string `xml:"index,attr,omitempty"`
	Initialization     string `xml:"initialization,attr,omitempty"`
	BitstreamSwitching string `xml:"bitstreamSwitching,attr,omitempty"`
}

type SegmentTimelineType struct {
	T     uint64   `xml:"t,attr,omitempty"`
	N     uint64   `xml:"n,attr,omitempty"`
	D     uint64   `xml:"d,attr"`
	R     int      `xml:"r,attr,omitempty"`
	Items []string `xml:",any"`
	S     []S      `xml:"urn:mpeg:dash:schema:mpd:2011 S"`
}

func (t *SegmentTimelineType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type T SegmentTimelineType
	var overlay struct {
		*T
		R int `xml:"r,attr,omitempty"`
	}
	overlay.T = (*T)(t)
	if err := d.DecodeElement(&overlay, &start); err != nil {
		return err
	}
	overlay.T.R = int(overlay.R)
	return nil
}

type SegmentURLType struct {
	Media      string   `xml:"media,attr,omitempty"`
	MediaRange string   `xml:"mediaRange,attr,omitempty"`
	Index      string   `xml:"index,attr,omitempty"`
	IndexRange string   `xml:"indexRange,attr,omitempty"`
	Items      []string `xml:",any"`
}

// Must match the pattern [^\r\n\t \p{Z}]*
type StringNoWhitespaceType string

type StringVectorType []string

func (x *StringVectorType) MarshalText() ([]byte, error) {
	result := make([][]byte, 0, len(*x))
	for _, v := range *x {
		result = append(result, []byte(v))
	}
	return bytes.Join(result, []byte(" ")), nil
}
func (x *StringVectorType) UnmarshalText(text []byte) error {
	for _, v := range bytes.Fields(text) {
		*x = append(*x, string(v))
	}
	return nil
}

type SubRepresentationType struct {
	RepresentationBaseType
	Level            uint             `xml:"level,attr,omitempty"`
	DependencyLevel  UIntVectorType   `xml:"dependencyLevel,attr,omitempty"`
	Bandwidth        uint             `xml:"bandwidth,attr,omitempty"`
	ContentComponent StringVectorType `xml:"contentComponent,attr,omitempty"`
}

type SubsetType struct {
	Contains UIntVectorType `xml:"contains,attr"`
}

type UIntVectorType []uint

func (x *UIntVectorType) MarshalText() ([]byte, error) {
	result := make([][]byte, 0, len(*x))
	for _, v := range *x {
		result = append(result, []byte(strconv.FormatUint(uint64(v), 10)))
	}
	return bytes.Join(result, []byte(" ")), nil
}
func (x *UIntVectorType) UnmarshalText(text []byte) error {
	for _, v := range strings.Fields(string(text)) {
		if i, err := strconv.ParseUint(v, 10, 32); err != nil {
			return err
		} else {
			*x = append(*x, uint(i))
		}
	}
	return nil
}

type URLType struct {
	SourceURL string   `xml:"sourceURL,attr,omitempty"`
	Range     string   `xml:"range,attr,omitempty"`
	Items     []string `xml:",any"`
}

// May be one of progressive, interlaced, unknown
type VideoScanType string

type xsdDateTime time.Time

func (t *xsdDateTime) UnmarshalText(text []byte) error {
	return _unmarshalTime(text, (*time.Time)(t), "2006-01-02T15:04:05.999999999")
}
func (t xsdDateTime) MarshalText() ([]byte, error) {
	return []byte((time.Time)(t).Format("2006-01-02T15:04:05.999999999")), nil
}
func _unmarshalTime(text []byte, t *time.Time, format string) (err error) {
	s := string(bytes.TrimSpace(text))
	*t, err = time.Parse(format, s)
	if _, ok := err.(*time.ParseError); ok {
		*t, err = time.Parse(format+"Z07:00", s)
	}
	return err
}
