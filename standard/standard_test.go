package standard_test

import (
	"testing"
	"time"

	"github.com/twpayne/go-nmea"
	"github.com/twpayne/go-nmea/nmeatest"
	"github.com/twpayne/go-nmea/standard"
)

func TestUblox(t *testing.T) {
	// From https://content.u-blox.com/sites/default/files/products/documents/u-blox8-M8_ReceiverDescrProtSpec_UBX-13003221.pdf.
	nmeatest.TestSentenceParserFunc(t,
		[]nmea.ParserOption{
			nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineStrict),
			nmea.WithLineEndingDiscipline(nmea.LineEndingDisciplineNever),
			nmea.WithSentenceParserFunc(standard.SentenceParserFunc),
		},
		[]nmeatest.TestCase{
			{
				S: "$GPDTM,W84,,0.0,N,0.0,E,0.0,W84*6F",
				Expected: &standard.DTM{
					Address:  nmea.NewAddress("GPDTM"),
					Datum:    "W84",
					RefDatum: "W84",
				},
			},
			{
				S: "$GPGBS,235503.00,1.6,1.4,3.2,,,,,,*40",
				Expected: &standard.GBS{
					Address: nmea.NewAddress("GPGBS"),
					TimeOfDay: nmea.TimeOfDay{
						Hour:   23,
						Minute: 55,
						Second: 3,
					},
					ErrLat: 1.6,
					ErrLon: 1.4,
					ErrAlt: 3.2,
				},
			},
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineRequire),
				},
				S: "$GPGBS,235458.00,1.4,1.3,3.1,03,,-21.4,3.8,1,0*5B",
				Expected: &standard.GBS{
					Address: nmea.NewAddress("GPGBS"),
					TimeOfDay: nmea.TimeOfDay{
						Hour:   23,
						Minute: 54,
						Second: 58,
					},
					ErrLat:   1.4,
					ErrLon:   1.3,
					ErrAlt:   3.1,
					SVID:     nmea.NewOptional(3),
					Bias:     nmea.NewOptional(-21.4),
					StdDev:   nmea.NewOptional(3.8),
					SystemID: nmea.NewOptional(1),
					SignalID: nmea.NewOptional(0),
				},
			},
			{
				S: "$GPGGA,092725.00,4717.11399,N,00833.91590,E,1,08,1.01,499.6,M,48.0,M,,*5B",
				Expected: &standard.GGA{
					Address: nmea.NewAddress("GPGGA"),
					TimeOfDay: nmea.TimeOfDay{
						Hour:   9,
						Minute: 27,
						Second: 25,
					},
					Lat:                              47.285233166666664,
					Lon:                              8.565265,
					FixQuality:                       1,
					NumberOfSatellites:               8,
					HDOP:                             1.01,
					Alt:                              499.6,
					HeightOfGeoidAboveWGS84Ellipsoid: 48,
				},
			},
			{
				S: "$GPGLL,4717.11364,N,00833.91565,E,092321.00,A,A*60",
				Expected: &standard.GLL{
					Address: nmea.NewAddress("GPGLL"),
					Lat:     47.28522733333333,
					Lon:     8.565260833333333,
					TimeOfDay: nmea.TimeOfDay{
						Hour:   9,
						Minute: 23,
						Second: 21,
					},
					Status:  'A',
					PosMode: 'A',
				},
			},
			{
				S: "$GNGNS,103600.01,5114.51176,N,00012.29380,W,ANNN,07,1.18,111.5,45.6,,,V*00",
				Expected: &standard.GNS{
					Address: nmea.NewAddress("GNGNS"),
					TimeOfDay: nmea.TimeOfDay{
						Hour:       10,
						Minute:     36,
						Nanosecond: 10000000,
					},
					Lat:       nmea.NewOptional(51.24186266666667),
					Lon:       nmea.NewOptional(-0.20489666666666664),
					PosMode:   []byte{'A', 'N', 'N', 'N'},
					NumSV:     7,
					HDOP:      nmea.NewOptional(1.18),
					Alt:       nmea.NewOptional(111.5),
					Sep:       nmea.NewOptional(45.6),
					NavStatus: 'V',
				},
			},
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineRequire),
				},
				S: "$GNGNS,122310.2,3722.425671,N,12258.856215,W,DAAA,14,0.9,1005.543,6.5,,,V*0E",
				Expected: &standard.GNS{
					Address: nmea.NewAddress("GNGNS"),
					TimeOfDay: nmea.TimeOfDay{
						Hour:       12,
						Minute:     23,
						Second:     10,
						Nanosecond: 200000000,
					},
					Lat:       nmea.NewOptional(37.373761183333336),
					Lon:       nmea.NewOptional(-122.98093691666666),
					PosMode:   []byte{'D', 'A', 'A', 'A'},
					NumSV:     14,
					HDOP:      nmea.NewOptional(0.9),
					Alt:       nmea.NewOptional(1005.543),
					Sep:       nmea.NewOptional(6.5),
					NavStatus: 'V',
				},
			},
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineRequire),
				},
				S: "$GPGNS,122310.2,,,,,,07,,,,5.2,23,V*02",
				Expected: &standard.GNS{
					Address: nmea.NewAddress("GPGNS"),
					TimeOfDay: nmea.TimeOfDay{
						Hour:       12,
						Minute:     23,
						Second:     10,
						Nanosecond: 200000000,
					},
					NumSV:       7,
					DiffAge:     nmea.NewOptional(5.2),
					DiffStation: nmea.NewOptional(23),
					NavStatus:   'V',
				},
			},
			{
				S: "$GNGRS,104148.00,1,2.6,2.2,-1.6,-1.1,-1.7,-1.5,5.8,1.7,,,,,1,1*52",
				Expected: &standard.GRS{
					Address: nmea.NewAddress("GNGRS"),
					TimeOfDay: nmea.TimeOfDay{
						Hour:   10,
						Minute: 41,
						Second: 48,
					},
					Mode: 1,
					Residuals: []nmea.Optional[float64]{
						nmea.NewOptional(2.6),
						nmea.NewOptional(2.2),
						nmea.NewOptional(-1.6),
						nmea.NewOptional(-1.1),
						nmea.NewOptional(-1.7),
						nmea.NewOptional(-1.5),
						nmea.NewOptional(5.8),
						nmea.NewOptional(1.7),
						{},
						{},
						{},
						{},
					},
					SystemID: 1,
					SignalID: 1,
				},
			},
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineRequire),
				},
				S: "$GNGRS,104148.00,1,,0.0,2.5,0.0,,2.8,,,,,,,1,5*52",
				Expected: &standard.GRS{
					Address: nmea.NewAddress("GNGRS"),
					TimeOfDay: nmea.TimeOfDay{
						Hour:   10,
						Minute: 41,
						Second: 48,
					},
					Mode: 1,
					Residuals: []nmea.Optional[float64]{
						{},
						nmea.NewOptional(0.0),
						nmea.NewOptional(2.5),
						nmea.NewOptional(0.0),
						{},
						nmea.NewOptional(2.8),
						{},
						{},
						{},
						{},
						{},
						{},
					},
					SystemID: 1,
					SignalID: 5,
				},
			},
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineRequire),
				},
				S: "$GPGSA,A,3,23,29,07,08,09,18,26,28,,,,,1.94,1.18,1.54,1*0D",
				Expected: &standard.GSA{
					Address: nmea.NewAddress("GPGSA"),
					OpMode:  'A',
					NavMode: 3,
					SVIDs: []nmea.Optional[int]{
						nmea.NewOptional(23),
						nmea.NewOptional(29),
						nmea.NewOptional(7),
						nmea.NewOptional(8),
						nmea.NewOptional(9),
						nmea.NewOptional(18),
						nmea.NewOptional(26),
						nmea.NewOptional(28),
						{},
						{},
						{},
						{},
					},
					PDOP:     nmea.NewOptional(1.94),
					HDOP:     nmea.NewOptional(1.18),
					VDOP:     nmea.NewOptional(1.54),
					SystemID: nmea.NewOptional(1),
				},
			},
			{
				S: "$GPGST,082356.00,1.8,,,,1.7,1.3,2.2*7E",
				Expected: &standard.GST{
					Address: nmea.NewAddress("GPGST"),
					TimeOfDay: nmea.TimeOfDay{
						Hour:   8,
						Minute: 23,
						Second: 56,
					},
					RangeRMS:  1.8,
					LatStdDev: 1.7,
					LonStdDev: 1.3,
					AltStdDev: 2.2,
				},
			},
			{
				S: "$GPGSV,3,1,09,09,,,17,10,,,40,12,,,49,13,,,35,1*6F",
				Expected: &standard.GSV{
					Address: nmea.NewAddress("GPGSV"),
					NumMsg:  3,
					MsgNum:  1,
					NumSV:   9,
					SatellitesInView: []standard.SatelliteInView{
						{
							SVID: 9,
							CNO:  nmea.NewOptional(17),
						},
						{
							SVID: 10,
							CNO:  nmea.NewOptional(40),
						},
						{
							SVID: 12,
							CNO:  nmea.NewOptional(49),
						},
						{
							SVID: 13,
							CNO:  nmea.NewOptional(35),
						},
					},
					SignalID: nmea.NewOptional(1),
				},
			},
			{
				S: "$GPGSV,3,2,09,15,,,44,17,,,45,19,,,44,24,,,50,1*64",
				Expected: &standard.GSV{
					Address: nmea.NewAddress("GPGSV"),
					NumMsg:  3,
					MsgNum:  2,
					NumSV:   9,
					SatellitesInView: []standard.SatelliteInView{
						{
							SVID: 15,
							CNO:  nmea.NewOptional(44),
						},
						{
							SVID: 17,
							CNO:  nmea.NewOptional(45),
						},
						{
							SVID: 19,
							CNO:  nmea.NewOptional(44),
						},
						{
							SVID: 24,
							CNO:  nmea.NewOptional(50),
						},
					},
					SignalID: nmea.NewOptional(1),
				},
			},
			{
				S: "$GPGSV,3,3,09,25,,,40,1*6E",
				Expected: &standard.GSV{
					Address: nmea.NewAddress("GPGSV"),
					NumMsg:  3,
					MsgNum:  3,
					NumSV:   9,
					SatellitesInView: []standard.SatelliteInView{
						{
							SVID: 25,
							CNO:  nmea.NewOptional(40),
						},
					},
					SignalID: nmea.NewOptional(1),
				},
			},
			{
				S: "$GPGSV,1,1,03,12,,,42,24,,,47,32,,,37,5*66",
				Expected: &standard.GSV{
					Address: nmea.NewAddress("GPGSV"),
					NumMsg:  1,
					MsgNum:  1,
					NumSV:   3,
					SatellitesInView: []standard.SatelliteInView{
						{
							SVID: 12,
							CNO:  nmea.NewOptional(42),
						},
						{
							SVID: 24,
							CNO:  nmea.NewOptional(47),
						},
						{
							SVID: 32,
							CNO:  nmea.NewOptional(37),
						},
					},
					SignalID: nmea.NewOptional(5),
				},
			},
			{
				S: "$GAGSV,1,1,00,2*76",
				Expected: &standard.GSV{
					Address:  nmea.NewAddress("GAGSV"),
					NumMsg:   1,
					MsgNum:   1,
					NumSV:    0,
					SignalID: nmea.NewOptional(2),
				},
			},
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineRequire),
				},
				S: "$GPRMC,083559.00,A,4717.11437,N,00833.91522,E,0.004,77.52,091202,,,A,V*57",
				Expected: &standard.RMC{
					Address:           nmea.NewAddress("GPRMC"),
					Time:              time.Date(2002, time.December, 9, 8, 35, 59, 0, time.UTC),
					Status:            'A',
					Lat:               47.2852395,
					Lon:               8.565253666666667,
					SpeedOverGroundKN: 0.004,
					CourseOverGround:  nmea.NewOptional(77.52),
					ModeIndicator:     'A',
					NavStatus:         'V',
				},
			},
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineRequire),
				},
				S: "$GPTHS,77.52,E*32",
				Expected: &standard.THS{
					Address:       nmea.NewAddress("GPTHS"),
					TrueHeading:   77.52,
					ModeIndicator: 'E',
				},
			},
			{
				S: "$GPTXT,01,01,02,u-blox ag - www.u-blox.com*50",
				Expected: &standard.TXT{
					Address: nmea.NewAddress("GPTXT"),
					NumMsg:  1,
					MsgNum:  1,
					MsgType: 2,
					Text:    "u-blox ag - www.u-blox.com",
				},
			},
			{
				S: "$GPTXT,01,01,02,ANTARIS ATR0620 HW 00000040*67",
				Expected: &standard.TXT{
					Address: nmea.NewAddress("GPTXT"),
					NumMsg:  1,
					MsgNum:  1,
					MsgType: 2,
					Text:    "ANTARIS ATR0620 HW 00000040",
				},
			},
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineRequire),
				},
				S: "$GPVLW,,N,,N,15.8,N,1.2,N*06",
				Expected: &standard.VLW{
					Address:               nmea.NewAddress("GPVLW"),
					TotalGroundDistanceNM: nmea.NewOptional(15.8),
					GroundDistanceNM:      nmea.NewOptional(1.2),
				},
			},
			{
				S: "$GPVTG,77.52,T,,M,0.004,N,0.008,K,A*06",
				Expected: &standard.VTG{
					Address:              nmea.NewAddress("GPVTG"),
					TrueCourseOverGround: 77.52,
					SpeedOverGroundKN:    0.004,
					SpeedOverGroundKPH:   0.008,
					ModeIndicator:        'A',
				},
			},
			{
				S: "$GPZDA,082710.00,16,09,2002,00,00*64",
				Expected: &standard.ZDA{
					Address:              nmea.NewAddress("GPZDA"),
					Time:                 time.Date(2002, time.September, 16, 8, 27, 10, 0, time.UTC),
					LocalTimeZoneHours:   0,
					LocalTimeZoneMinutes: 0,
				},
			},
		},
	)
}

func TestSparkfun(t *testing.T) {
	// From https://www.sparkfun.com/datasheets/GPS/NMEA%20Reference%20Manual-Rev2.1-Dec07.pdf.
	nmeatest.TestSentenceParserFunc(t,
		[]nmea.ParserOption{
			nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineStrict),
			nmea.WithLineEndingDiscipline(nmea.LineEndingDisciplineNever),
			nmea.WithSentenceParserFunc(standard.SentenceParserFunc),
		},
		[]nmeatest.TestCase{
			{
				S: "$GPGGA,002153.000,3342.6618,N,11751.3858,W,1,10,1.2,27.0,M,-34.2,M,,0000*5E",
				Expected: &standard.GGA{
					Address: nmea.NewAddress("GPGGA"),
					TimeOfDay: nmea.TimeOfDay{
						Hour:   0,
						Minute: 21,
						Second: 53,
					},
					Lat:                              33.71103,
					Lon:                              -117.85643,
					FixQuality:                       1,
					NumberOfSatellites:               10,
					HDOP:                             1.2,
					Alt:                              27,
					HeightOfGeoidAboveWGS84Ellipsoid: -34.2,
					DGPSReferenceStationID:           "0000",
				},
			},
			{
				S: "$GPGLL,3723.2475,N,12158.3416,W,161229.487,A,A*41",
				Expected: &standard.GLL{
					Address: nmea.NewAddress("GPGLL"),
					Lat:     37.387458333333335,
					Lon:     -121.97236,
					TimeOfDay: nmea.TimeOfDay{
						Hour:       16,
						Minute:     12,
						Second:     29,
						Nanosecond: 487000000,
					},
					Status:  'A',
					PosMode: 'A',
				},
			},
			{
				S: "$GPGSA,A,3,07,02,26,27,09,04,15,,,,,,1.8,1.0,1.5*33",
				Expected: &standard.GSA{
					Address: nmea.NewAddress("GPGSA"),
					OpMode:  'A',
					NavMode: 3,
					SVIDs: []nmea.Optional[int]{
						nmea.NewOptional(7),
						nmea.NewOptional(2),
						nmea.NewOptional(26),
						nmea.NewOptional(27),
						nmea.NewOptional(9),
						nmea.NewOptional(4),
						nmea.NewOptional(15),
						{},
						{},
						{},
						{},
						{},
					},
					PDOP: nmea.NewOptional(1.8),
					HDOP: nmea.NewOptional(1.0),
					VDOP: nmea.NewOptional(1.5),
				},
			},
			{
				S: "$GPGSV,2,1,07,07,79,048,42,02,51,062,43,26,36,256,42,27,27,138,42*71",
				Expected: &standard.GSV{
					Address: nmea.NewAddress("GPGSV"),
					NumMsg:  2,
					MsgNum:  1,
					NumSV:   7,
					SatellitesInView: []standard.SatelliteInView{
						{
							SVID: 7,
							Elv:  nmea.NewOptional(79),
							Az:   nmea.NewOptional(48),
							CNO:  nmea.NewOptional(42),
						},
						{
							SVID: 2,
							Elv:  nmea.NewOptional(51),
							Az:   nmea.NewOptional(62),
							CNO:  nmea.NewOptional(43),
						},
						{
							SVID: 26,
							Elv:  nmea.NewOptional(36),
							Az:   nmea.NewOptional(256),
							CNO:  nmea.NewOptional(42),
						},
						{
							SVID: 27,
							Elv:  nmea.NewOptional(27),
							Az:   nmea.NewOptional(138),
							CNO:  nmea.NewOptional(42),
						},
					},
				},
			},
			{
				S: "$GPGSV,2,2,07,09,23,313,42,04,19,159,41,15,12,041,42*41",
				Expected: &standard.GSV{
					Address: nmea.NewAddress("GPGSV"),
					NumMsg:  2,
					MsgNum:  2,
					NumSV:   7,
					SatellitesInView: []standard.SatelliteInView{
						{
							SVID: 9,
							Elv:  nmea.NewOptional(23),
							Az:   nmea.NewOptional(313),
							CNO:  nmea.NewOptional(42),
						},
						{
							SVID: 4,
							Elv:  nmea.NewOptional(19),
							Az:   nmea.NewOptional(159),
							CNO:  nmea.NewOptional(41),
						},
						{
							SVID: 15,
							Elv:  nmea.NewOptional(12),
							Az:   nmea.NewOptional(41),
							CNO:  nmea.NewOptional(42),
						},
					},
				},
			},
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineRequire),
				},
				S: "$GPMSS,55,27,318.0,100,1*57",
				Expected: &standard.MSS{
					Address:            nmea.NewAddress("GPMSS"),
					SignalStrength:     55,
					SignalToNoiseRatio: 27,
					BeaconFrequencyKHz: 318,
					BeaconBitRate:      100,
					ChannelNumber:      nmea.NewOptional(1),
				},
			},
			{
				S: "$GPRMC,161229.487,A,3723.2475,N,12158.3416,W,0.13,309.62,120598,,*10",
				Expected: &standard.RMC{
					Address:           nmea.NewAddress("GPRMC"),
					Time:              time.Date(1998, time.May, 12, 16, 12, 29, 487000000, time.UTC),
					Status:            65,
					Lat:               37.387458333333335,
					Lon:               -121.97236,
					SpeedOverGroundKN: 0.13,
					CourseOverGround:  nmea.NewOptional(309.62),
				},
			},
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineRequire),
				},
				S: "$GPVTG,309.62,T,,M,0.13,N,0.2,K,A*23",
				Expected: &standard.VTG{
					Address:              nmea.NewAddress("GPVTG"),
					TrueCourseOverGround: 309.62,
					SpeedOverGroundKN:    0.13,
					SpeedOverGroundKPH:   0.2,
					ModeIndicator:        'A',
				},
			},
			{
				S: "$GPZDA,181813,14,10,2003,00,00*4F",
				Expected: &standard.ZDA{
					Address: nmea.NewAddress("GPZDA"),
					Time:    time.Date(2003, time.October, 14, 18, 18, 13, 0, time.UTC),
				},
			},
		},
	)
}

func TestGNSSDO(t *testing.T) {
	// From https://ww1.microchip.com/downloads/aemDocuments/documents/VOP/ProductDocuments/ReferenceManuals/GNSSDO_NMEA_Reference_Manual_RevA.pdf.
	nmeatest.TestSentenceParserFunc(t,
		[]nmea.ParserOption{
			nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineStrict),
			nmea.WithLineEndingDiscipline(nmea.LineEndingDisciplineNever),
			nmea.WithSentenceParserFunc(standard.SentenceParserFunc),
		},
		[]nmeatest.TestCase{
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineIgnore),
				},
				S: "$GNGSA,A,3,09,15,26,05,24,21,08,02,29,28,18,10,0.8,0.5,0.5,1*XX",
				Expected: &standard.GSA{
					Address: nmea.NewAddress("GNGSA"),
					OpMode:  'A',
					NavMode: 3,
					SVIDs: []nmea.Optional[int]{
						nmea.NewOptional(9),
						nmea.NewOptional(15),
						nmea.NewOptional(26),
						nmea.NewOptional(5),
						nmea.NewOptional(24),
						nmea.NewOptional(21),
						nmea.NewOptional(8),
						nmea.NewOptional(2),
						nmea.NewOptional(29),
						nmea.NewOptional(28),
						nmea.NewOptional(18),
						nmea.NewOptional(10),
					},
					PDOP:     nmea.NewOptional(0.8),
					HDOP:     nmea.NewOptional(0.5),
					VDOP:     nmea.NewOptional(0.5),
					SystemID: nmea.NewOptional(1),
				},
			},
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineIgnore),
				},
				S: "$GPGSV,4,1,14,15,67,319,52,09,63,068,53,26,45,039,50,05,44,104,49,1*XX",
				Expected: &standard.GSV{
					Address: nmea.NewAddress("GPGSV"),
					NumMsg:  4,
					MsgNum:  1,
					NumSV:   14,
					SatellitesInView: []standard.SatelliteInView{
						{
							SVID: 15,
							Elv:  nmea.NewOptional(67),
							Az:   nmea.NewOptional(319),
							CNO:  nmea.NewOptional(52),
						},
						{
							SVID: 9,
							Elv:  nmea.NewOptional(63),
							Az:   nmea.NewOptional(68),
							CNO:  nmea.NewOptional(53),
						},
						{
							SVID: 26,
							Elv:  nmea.NewOptional(45),
							Az:   nmea.NewOptional(39),
							CNO:  nmea.NewOptional(50),
						},
						{
							SVID: 5,
							Elv:  nmea.NewOptional(44),
							Az:   nmea.NewOptional(104),
							CNO:  nmea.NewOptional(49),
						},
					},
					SignalID: nmea.NewOptional(1),
				},
			},
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineIgnore),
				},
				S: "$GPGSV,4,2,14,24,42,196,47,21,34,302,46,18,12,305,43,28,11,067,41,1*XX",
				Expected: &standard.GSV{
					Address: nmea.NewAddress("GPGSV"),
					NumMsg:  4,
					MsgNum:  2,
					NumSV:   14,
					SatellitesInView: []standard.SatelliteInView{
						{
							SVID: 24,
							Elv:  nmea.NewOptional(42),
							Az:   nmea.NewOptional(196),
							CNO:  nmea.NewOptional(47),
						},
						{
							SVID: 21,
							Elv:  nmea.NewOptional(34),
							Az:   nmea.NewOptional(302),
							CNO:  nmea.NewOptional(46),
						},
						{
							SVID: 18,
							Elv:  nmea.NewOptional(12),
							Az:   nmea.NewOptional(305),
							CNO:  nmea.NewOptional(43),
						},
						{
							SVID: 28,
							Elv:  nmea.NewOptional(11),
							Az:   nmea.NewOptional(67),
							CNO:  nmea.NewOptional(41),
						},
					},
					SignalID: nmea.NewOptional(1),
				},
			},
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineIgnore),
				},
				S: "$GPGSV,4,3,14,08,07,035,38,29,04,237,39,02,02,161,40,50,47,163,44,1*XX",
				Expected: &standard.GSV{
					Address: nmea.NewAddress("GPGSV"),
					NumMsg:  4,
					MsgNum:  3,
					NumSV:   14,
					SatellitesInView: []standard.SatelliteInView{
						{
							SVID: 8,
							Elv:  nmea.NewOptional(7),
							Az:   nmea.NewOptional(35),
							CNO:  nmea.NewOptional(38),
						},
						{
							SVID: 29,
							Elv:  nmea.NewOptional(4),
							Az:   nmea.NewOptional(237),
							CNO:  nmea.NewOptional(39),
						},
						{
							SVID: 2,
							Elv:  nmea.NewOptional(2),
							Az:   nmea.NewOptional(161),
							CNO:  nmea.NewOptional(40),
						},
						{
							SVID: 50,
							Elv:  nmea.NewOptional(47),
							Az:   nmea.NewOptional(163),
							CNO:  nmea.NewOptional(44),
						},
					},
					SignalID: nmea.NewOptional(1),
				},
			},
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineIgnore),
				},
				S: "$GLGSV,3,1,09,79,66,099,50,69,55,019,53,80,33,176,46,68,28,088,45,1*XX",
				Expected: &standard.GSV{
					Address: nmea.NewAddress("GLGSV"),
					NumMsg:  3,
					MsgNum:  1,
					NumSV:   9,
					SatellitesInView: []standard.SatelliteInView{
						{
							SVID: 79,
							Elv:  nmea.NewOptional(66),
							Az:   nmea.NewOptional(99),
							CNO:  nmea.NewOptional(50),
						},
						{
							SVID: 69,
							Elv:  nmea.NewOptional(55),
							Az:   nmea.NewOptional(19),
							CNO:  nmea.NewOptional(53),
						},
						{
							SVID: 80,
							Elv:  nmea.NewOptional(33),
							Az:   nmea.NewOptional(176),
							CNO:  nmea.NewOptional(46),
						},
						{
							SVID: 68,
							Elv:  nmea.NewOptional(28),
							Az:   nmea.NewOptional(88),
							CNO:  nmea.NewOptional(45),
						},
					},
					SignalID: nmea.NewOptional(1),
				},
			},
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineIgnore),
				},
				S: "$GLGSV,3,2,09,70,25,315,46,78,24,031,42,85,18,293,44,84,16,246,41,1*XX",
				Expected: &standard.GSV{
					Address: nmea.NewAddress("GLGSV"),
					NumMsg:  3,
					MsgNum:  2,
					NumSV:   9,
					SatellitesInView: []standard.SatelliteInView{
						{
							SVID: 70,
							Elv:  nmea.NewOptional(25),
							Az:   nmea.NewOptional(315),
							CNO:  nmea.NewOptional(46),
						},
						{
							SVID: 78,
							Elv:  nmea.NewOptional(24),
							Az:   nmea.NewOptional(31),
							CNO:  nmea.NewOptional(42),
						},
						{
							SVID: 85,
							Elv:  nmea.NewOptional(18),
							Az:   nmea.NewOptional(293),
							CNO:  nmea.NewOptional(44),
						},
						{
							SVID: 84,
							Elv:  nmea.NewOptional(16),
							Az:   nmea.NewOptional(246),
							CNO:  nmea.NewOptional(41),
						},
					},
					SignalID: nmea.NewOptional(1),
				},
			},
			{
				Options: []nmea.ParserOption{
					nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineIgnore),
				},
				S: "$GNVTG,0.00,T,,M,0.00,N,0.00,K,D*XX",
				Expected: &standard.VTG{
					Address:       nmea.NewAddress("GNVTG"),
					ModeIndicator: 'D',
				},
			},
		},
	)
}

func TestMiscellaneous(t *testing.T) {
	nmeatest.TestSentenceParserFunc(t,
		[]nmea.ParserOption{
			nmea.WithChecksumDiscipline(nmea.ChecksumDisciplineStrict),
			nmea.WithLineEndingDiscipline(nmea.LineEndingDisciplineNever),
			nmea.WithSentenceParserFunc(standard.SentenceParserFunc),
		},
		[]nmeatest.TestCase{
			{
				S: "$GPMSS,0,0,0.000000,0,*58",
				Expected: &standard.MSS{
					Address: nmea.NewAddress("GPMSS"),
				},
			},
			{
				S: "$GPMSS,0,0,0.000000,200,*5A",
				Expected: &standard.MSS{
					Address:       nmea.NewAddress("GPMSS"),
					BeaconBitRate: 200,
				},
			},
			{
				S: "$GPMSS,55,27,318.0,100,*66",
				Expected: &standard.MSS{
					Address:            nmea.NewAddress("GPMSS"),
					SignalStrength:     55,
					SignalToNoiseRatio: 27,
					BeaconFrequencyKHz: 318,
					BeaconBitRate:      100,
				},
			},
			{
				Skip: "FIXME parse missing data",
				S:    "$GPRMC,102042.00,V,,,,,,,110324,,,N*7D",
			},
			{
				Skip: "FIXME parse missing data",
				S:    "$GPGGA,102039.00,,,,,0,00,99.99,,,,,,*6F",
			},
		},
	)
}
