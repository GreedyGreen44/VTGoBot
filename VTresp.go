package main

import "fmt"

type VtResp struct {
	AisDataset []VesselInfo `json:"AIS"`
}

type VesselInfo struct {
	Mmsi               int64   `json:"MMSI"`
	Timestamp          string  `json:"TIMESTAMP"`
	Latitude           float32 `json:"LATITUDE"`
	Longitude          float32 `json:"LONGITUDE"`
	Course             float32 `json:"COURSE"`
	Speed              float32 `json:"SPEED"`
	Heading            int32   `json:"HEADING"`
	Navstat            int32   `json:"NAVSTAT"`
	Imo                int64   `json:"IMO"`
	Name               string  `json:"NAME"`
	Callsign           string  `json:"CALLSIGN"`
	Type               int32   `json:"TYPE"`
	A                  int32   `json:"A"`
	B                  int32   `json:"B"`
	C                  int32   `json:"C"`
	D                  int32   `json:"D"`
	Draught            float32 `json:"DRAUGHT"`
	Destination        string  `json:"DESTINATION"`
	Locode             string  `json:"LOCODE"`
	ETA_Ais            string  `json:"ETA_AIS"`
	ETA                string  `json:"ETA"`
	ETA_predicted      string  `json:"ETA_PREDICTED"`
	Distance_remaining int32   `json:"DISTANCE_REMAINING"`
	Src                string  `json:"SRC"`
	Zone               string  `json:"ZONE"`
	Eca                bool    `json:"ECA"`
}

func decodeNavSat(navSat int32) (str string) {
	switch navSat {
	case 0:
		return "Under way using engine"
	case 1:
		return "At anchor"
	case 2:
		return "Not under command"
	case 3:
		return "Restricted manoeuverability"
	case 4:
		return "Constrained by her draught"
	case 5:
		return "Moored"
	case 6:
		return "Aground"
	case 7:
		return "Engaged in Fishing"
	case 8:
		return "Under way sailing"
	case 9:
		return "Reserved for future amendment of Navigational Status for HSC"
	case 10:
		return "Reserved for future amendment of Navigational Status for WIG"
	case 11:
		return "Reserved for future use"
	case 12:
		return "Reserved for future use"
	case 13:
		return "Reserved for future use"
	case 14:
		return "AIS-SART is active"
	default:
		return "Not defined (default)"
	}
}

func decodeShipType(typeCode int32) (descr string) {
	code := int(typeCode)

	switch {
	case code == 0:
		return "Not available (default)"
	case code >= 1 && code <= 19:
		return "Reserved for future use"
	case code == 20:
		return "Wing in ground (WIG), all ships of this type"
	case code >= 21 && code <= 24:
		return fmt.Sprintf("Wing in ground (WIG), Hazardous category %c", 'A'+(code-21))
	case code >= 25 && code <= 29:
		return "Wing in ground (WIG), Reserved for future use"
	case code == 30:
		return "Fishing"
	case code == 31:
		return "Towing"
	case code == 32:
		return "Towing: length exceeds 200m or breadth exceeds 25m"
	case code == 33:
		return "Dredging or underwater ops"
	case code == 34:
		return "Diving ops"
	case code == 35:
		return "Military ops"
	case code == 36:
		return "Sailing"
	case code == 37:
		return "Pleasure Craft"
	case code == 38 || code == 39:
		return "Reserved"
	case code == 40:
		return "High speed craft (HSC), all ships of this type"
	case code >= 41 && code <= 44:
		return fmt.Sprintf("High speed craft (HSC), Hazardous category %c", 'A'+(code-41))
	case code >= 45 && code <= 48:
		return "High speed craft (HSC), Reserved for future use"
	case code == 49:
		return "High speed craft (HSC), No additional information"
	case code == 50:
		return "Pilot Vessel"
	case code == 51:
		return "Search and Rescue vessel"
	case code == 52:
		return "Tug"
	case code == 53:
		return "Port Tender"
	case code == 54:
		return "Anti-pollution equipment"
	case code == 55:
		return "Law Enforcement"
	case code == 56 || code == 57:
		return "Spare - Local Vessel"
	case code == 58:
		return "Medical Transport"
	case code == 59:
		return "Noncombatant ship according to RR Resolution No. 18"
	case code == 60:
		return "Passenger, all ships of this type"
	case code >= 61 && code <= 64:
		return fmt.Sprintf("Passenger, Hazardous category %c", 'A'+(code-61))
	case code >= 65 && code <= 68:
		return "Passenger, Reserved for future use"
	case code == 69:
		return "Passenger, No additional information"
	case code == 70:
		return "Cargo, all ships of this type"
	case code >= 71 && code <= 74:
		return fmt.Sprintf("Cargo, Hazardous category %c", 'A'+(code-71))
	case code >= 75 && code <= 78:
		return "Cargo, Reserved for future use"
	case code == 79:
		return "Cargo, No additional information"
	case code == 80:
		return "Tanker, all ships of this type"
	case code >= 81 && code <= 84:
		return fmt.Sprintf("Tanker, Hazardous category %c", 'A'+(code-81))
	case code >= 85 && code <= 88:
		return "Tanker, Reserved for future use"
	case code == 89:
		return "Tanker, No additional information"
	case code == 90:
		return "Other Type, all ships of this type"
	case code >= 91 && code <= 94:
		return fmt.Sprintf("Other Type, Hazardous category %c", 'A'+(code-91))
	case code >= 95 && code <= 98:
		return "Other Type, Reserved for future use"
	case code == 99:
		return "Other Type, no additional information"
	default:
		return ""
	}
}
