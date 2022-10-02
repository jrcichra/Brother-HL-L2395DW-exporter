package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type BrotherCollector struct {
	url                                string
	timeout                            time.Duration
	memorySizeMetric                   *prometheus.Desc
	pageCounterMetric                  *prometheus.Desc
	averageCoverageMetric              *prometheus.Desc
	drumUnitPercentLifeRemainingMetric *prometheus.Desc
	TonerPercentLifeRemainingMetric    *prometheus.Desc
	A4LetterMetric                     *prometheus.Desc
	LegalFolioMetric                   *prometheus.Desc
	B5ExecutiveMetric                  *prometheus.Desc
	EnvelopesMetric                    *prometheus.Desc
	A5Metric                           *prometheus.Desc
	Others01Metric                     *prometheus.Desc
	PlainThinRecycledMetric            *prometheus.Desc
	ThickThickerBondMetric             *prometheus.Desc
	EnvelopesEnvThicEnvThinMetric      *prometheus.Desc
	LabelMetric                        *prometheus.Desc
	HagakiMetric                       *prometheus.Desc
	TotalMetric                        *prometheus.Desc
	Total2SidedPrintMetric             *prometheus.Desc
	CopyMetric                         *prometheus.Desc
	Copy2SidedPrintMetric              *prometheus.Desc
	PrintMetric                        *prometheus.Desc
	Print2SidedPrintMetric             *prometheus.Desc
	Others02Metric                     *prometheus.Desc
	Others2SidedPrintMetric            *prometheus.Desc
	FlatbedScanMetric                  *prometheus.Desc
	ScanPageCountMetric                *prometheus.Desc
	TonerReplaceCountMetric            *prometheus.Desc
	DrumReplaceCountMetric             *prometheus.Desc
	TotalPaperJamsMetric               *prometheus.Desc
	JamTray1Metric                     *prometheus.Desc
	JamInsideMetric                    *prometheus.Desc
	JamRearMetric                      *prometheus.Desc
	Jam2SidedMetric                    *prometheus.Desc
	ErrorCountMetric                   *prometheus.Desc
}

var labels = []string{
	"nodeName", "modelName", "location", "contact", "ipAddress",
	"serialNumber", "mainFirmwareVersion", "sub1FirmwareVersion",
}

func newBrotherCollector(url string, timeout time.Duration) *BrotherCollector {
	errorLabels := make([]string, len(labels))
	copy(errorLabels, labels)
	errorLabels = append(errorLabels, "error_message")
	return &BrotherCollector{
		url:                                url,
		timeout:                            timeout,
		memorySizeMetric:                   prometheus.NewDesc("brother_memory_size", "", labels, nil),
		pageCounterMetric:                  prometheus.NewDesc("brother_page_counter", "", labels, nil),
		averageCoverageMetric:              prometheus.NewDesc("brother_average_coverage", "", labels, nil),
		drumUnitPercentLifeRemainingMetric: prometheus.NewDesc("brother_drum_unit_percent_life_remaining", "", labels, nil),
		TonerPercentLifeRemainingMetric:    prometheus.NewDesc("brother_toner_percent_life_remaining", "", labels, nil),
		A4LetterMetric:                     prometheus.NewDesc("brother_page_counter_a4_letter", "", labels, nil),
		LegalFolioMetric:                   prometheus.NewDesc("brother_page_counter_legal_folio", "", labels, nil),
		B5ExecutiveMetric:                  prometheus.NewDesc("brother_page_counter_b5_executive", "", labels, nil),
		EnvelopesMetric:                    prometheus.NewDesc("brother_page_counter_envelopes", "", labels, nil),
		A5Metric:                           prometheus.NewDesc("brother_page_counter_a5", "", labels, nil),
		Others01Metric:                     prometheus.NewDesc("brother_others_01", "", labels, nil),
		PlainThinRecycledMetric:            prometheus.NewDesc("brother_page_counter_plain_thin_recycled", "", labels, nil),
		ThickThickerBondMetric:             prometheus.NewDesc("brother_page_counter_thick_thicker_bond", "", labels, nil),
		EnvelopesEnvThicEnvThinMetric:      prometheus.NewDesc("brother_page_counter_envelopes_env_thic_env_thin", "", labels, nil),
		LabelMetric:                        prometheus.NewDesc("brother_page_counter_label", "", labels, nil),
		HagakiMetric:                       prometheus.NewDesc("brother_page_counter_hagaki", "", labels, nil),
		TotalMetric:                        prometheus.NewDesc("brother_page_counter_total", "", labels, nil),
		Total2SidedPrintMetric:             prometheus.NewDesc("brother_page_counter_total_two_sided", "", labels, nil),
		CopyMetric:                         prometheus.NewDesc("brother_copies", "", labels, nil),
		Copy2SidedPrintMetric:              prometheus.NewDesc("brother_copies_two_sided", "", labels, nil),
		PrintMetric:                        prometheus.NewDesc("brother_prints", "", labels, nil),
		Print2SidedPrintMetric:             prometheus.NewDesc("brother_prints_two_sided", "", labels, nil),
		Others02Metric:                     prometheus.NewDesc("brother_others_02", "", labels, nil),
		Others2SidedPrintMetric:            prometheus.NewDesc("brother_others_two_sided", "", labels, nil),
		FlatbedScanMetric:                  prometheus.NewDesc("brother_scan", "", labels, nil),
		ScanPageCountMetric:                prometheus.NewDesc("brother_scan_page_counter", "", labels, nil),
		TonerReplaceCountMetric:            prometheus.NewDesc("brother_toner_replacements", "", labels, nil),
		DrumReplaceCountMetric:             prometheus.NewDesc("brother_drum_replacements", "", labels, nil),
		TotalPaperJamsMetric:               prometheus.NewDesc("brother_paper_jams", "", labels, nil),
		JamTray1Metric:                     prometheus.NewDesc("brother_paper_jam_tray_1", "", labels, nil),
		JamInsideMetric:                    prometheus.NewDesc("brother_paper_jam_inside", "", labels, nil),
		JamRearMetric:                      prometheus.NewDesc("brother_paper_jam_rear", "", labels, nil),
		Jam2SidedMetric:                    prometheus.NewDesc("brother_paper_jam_two_sided", "", labels, nil),
		ErrorCountMetric:                   prometheus.NewDesc("brother_error_count", "", errorLabels, nil),
	}
}

var _ prometheus.Collector = &BrotherCollector{}

func (b *BrotherCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- b.memorySizeMetric
	ch <- b.pageCounterMetric
	ch <- b.averageCoverageMetric
	ch <- b.drumUnitPercentLifeRemainingMetric
	ch <- b.TonerPercentLifeRemainingMetric
	ch <- b.A4LetterMetric
	ch <- b.LegalFolioMetric
	ch <- b.B5ExecutiveMetric
	ch <- b.EnvelopesMetric
	ch <- b.A5Metric
	ch <- b.Others01Metric
	ch <- b.PlainThinRecycledMetric
	ch <- b.ThickThickerBondMetric
	ch <- b.EnvelopesEnvThicEnvThinMetric
	ch <- b.LabelMetric
	ch <- b.HagakiMetric
	ch <- b.TotalMetric
	ch <- b.Total2SidedPrintMetric
	ch <- b.CopyMetric
	ch <- b.Copy2SidedPrintMetric
	ch <- b.PrintMetric
	ch <- b.Print2SidedPrintMetric
	ch <- b.Others02Metric
	ch <- b.Others2SidedPrintMetric
	ch <- b.FlatbedScanMetric
	ch <- b.ScanPageCountMetric
	ch <- b.TonerReplaceCountMetric
	ch <- b.DrumReplaceCountMetric
	ch <- b.TotalPaperJamsMetric
	ch <- b.JamTray1Metric
	ch <- b.JamInsideMetric
	ch <- b.JamRearMetric
	ch <- b.Jam2SidedMetric
	ch <- b.ErrorCountMetric
}

func (b *BrotherCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", b.url, nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != http.StatusOK {
		panic(fmt.Errorf("did not get a 200 OK: %d", res.StatusCode))
	}
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(res.Body)
	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}
	if len(records) <= 1 {
		panic(fmt.Errorf("no printer rows found. found %d rows in total", len(records)))
	}

	// should only be one record
	row := records[1]

	// convert what we can to floats
	values := make([]float64, len(row))
	for i, value := range row {
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			values[i] = f
		}
	}

	labelValues := row[0:8]
	labelValuesWithError := make([]string, len(labelValues))
	copy(labelValuesWithError, labelValues)
	labelValuesWithError = append(labelValuesWithError, "CHANGEME")
	labelValuesWithErrorLenMinusOne := len(labelValuesWithError) - 1

	// place the floats as values
	ch <- prometheus.MustNewConstMetric(b.memorySizeMetric, prometheus.GaugeValue, values[8], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.pageCounterMetric, prometheus.GaugeValue, values[9], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.averageCoverageMetric, prometheus.GaugeValue, values[10], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.drumUnitPercentLifeRemainingMetric, prometheus.GaugeValue, values[11], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.TonerPercentLifeRemainingMetric, prometheus.GaugeValue, values[12], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.A4LetterMetric, prometheus.GaugeValue, values[13], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.LegalFolioMetric, prometheus.GaugeValue, values[14], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.B5ExecutiveMetric, prometheus.GaugeValue, values[15], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.EnvelopesMetric, prometheus.GaugeValue, values[16], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.A5Metric, prometheus.GaugeValue, values[17], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.Others01Metric, prometheus.GaugeValue, values[18], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.PlainThinRecycledMetric, prometheus.GaugeValue, values[19], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.ThickThickerBondMetric, prometheus.GaugeValue, values[20], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.EnvelopesEnvThicEnvThinMetric, prometheus.GaugeValue, values[21], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.LabelMetric, prometheus.GaugeValue, values[22], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.HagakiMetric, prometheus.GaugeValue, values[23], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.TotalMetric, prometheus.GaugeValue, values[24], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.Total2SidedPrintMetric, prometheus.GaugeValue, values[25], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.CopyMetric, prometheus.GaugeValue, values[26], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.Copy2SidedPrintMetric, prometheus.GaugeValue, values[27], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.PrintMetric, prometheus.GaugeValue, values[28], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.Print2SidedPrintMetric, prometheus.GaugeValue, values[29], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.Others02Metric, prometheus.GaugeValue, values[30], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.Others2SidedPrintMetric, prometheus.GaugeValue, values[31], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.FlatbedScanMetric, prometheus.GaugeValue, values[32], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.ScanPageCountMetric, prometheus.GaugeValue, values[33], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.TonerReplaceCountMetric, prometheus.GaugeValue, values[34], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.DrumReplaceCountMetric, prometheus.GaugeValue, values[35], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.TotalPaperJamsMetric, prometheus.GaugeValue, values[36], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.JamTray1Metric, prometheus.GaugeValue, values[37], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.JamInsideMetric, prometheus.GaugeValue, values[38], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.JamRearMetric, prometheus.GaugeValue, values[39], labelValues...)
	ch <- prometheus.MustNewConstMetric(b.Jam2SidedMetric, prometheus.GaugeValue, values[40], labelValues...)

	// unique errors
	m := make(map[string]float64)
	for i := 41; i <= 50; i++ {
		m[row[i]] = values[i+10]
	}
	for key, val := range m {
		labelValuesWithError[labelValuesWithErrorLenMinusOne] = key
		ch <- prometheus.MustNewConstMetric(b.ErrorCountMetric, prometheus.GaugeValue, val, labelValuesWithError...)
	}
}

func main() {

	address := flag.String("address", "10.0.0.3", "IP address of the brother printer")
	csvURL := flag.String("csvURL", "etc/mnt_info.csv", "Path for the csv file on the printer")
	timeout := flag.Int("timeout", 10, "context timeout for HTTP call")
	flag.Parse()

	fullURL := fmt.Sprintf("http://%s/%s", *address, *csvURL)
	b := newBrotherCollector(fullURL, time.Duration(*timeout))
	prometheus.MustRegister(b)

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Beginning to serve on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
