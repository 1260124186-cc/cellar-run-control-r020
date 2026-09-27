package domain

import "math"

type RunMetrics struct {
	FinalGravity        float64 `json:"final_gravity"`
	ApparentAttenuation float64 `json:"apparent_attenuation_percent"`
	EstimatedAlcoholABV float64 `json:"estimated_alcohol_abv"`
	ObservationCount    int     `json:"observation_count"`
}

type ObservationWindow struct {
	observations []GravityObservation
}

func NewObservationWindow(
	observations []GravityObservation,
) ObservationWindow {
	return ObservationWindow{
		observations: append([]GravityObservation(nil), observations...),
	}
}

func (w ObservationWindow) ObservationCount() int {
	return len(w.observations)
}

func (w ObservationWindow) FinalObservation() (GravityObservation, bool) {
	if len(w.observations) == 0 {
		return GravityObservation{}, false
	}
	return w.observations[len(w.observations)-1], true
}

func (w ObservationWindow) Calculate(originalGravity float64) RunMetrics {
	final, ok := w.FinalObservation()
	if !ok {
		return RunMetrics{}
	}
	return CalculateMetrics(
		originalGravity,
		final.Gravity,
		w.ObservationCount(),
	)
}

func CalculateMetrics(originalGravity, finalGravity float64, count int) RunMetrics {
	attenuation := 0.0
	if originalGravity > 1 {
		attenuation = ((originalGravity - finalGravity) / (originalGravity - 1)) * 100
	}
	alcohol := (originalGravity - finalGravity) * 131.25
	if alcohol < 0 {
		alcohol = 0
	}
	return RunMetrics{
		FinalGravity:        finalGravity,
		ApparentAttenuation: round2(attenuation),
		EstimatedAlcoholABV: round2(alcohol),
		ObservationCount:    count,
	}
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}
