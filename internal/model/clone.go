package model

func CloneLayers(in []Layer) []Layer {
	if in == nil {
		return nil
	}
	out := make([]Layer, len(in))
	copy(out, in)
	return out
}

func CloneMeasurements(in []Measurement) []Measurement {
	if in == nil {
		return nil
	}
	out := make([]Measurement, len(in))
	copy(out, in)
	return out
}

func CloneAlerts(in []Alert) []Alert {
	if in == nil {
		return nil
	}
	out := make([]Alert, len(in))
	copy(out, in)
	return out
}

func CloneEvents(in []Event) []Event {
	if in == nil {
		return nil
	}
	out := make([]Event, len(in))
	copy(out, in)
	return out
}
