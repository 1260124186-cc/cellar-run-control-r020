package domain

import "math"

type RunMetrics struct {
	FinalGravity        float64 `json:"final_gravity"`
	ApparentAttenuation float64 `json:"apparent_attenuation_percent"`
	EstimatedAlcoholABV float64 `json:"estimated_alcohol_abv"`
	ObservationCount    int     `json:"observation_count"`
}

type ObservationWindowPolicy struct {
	TrimTrailingRise  bool
	MaximumStableRise float64
	MinimumSamples    int
}

type ObservationWindow struct {
	observations []GravityObservation
	policy       ObservationWindowPolicy
}

func NewObservationWindow(
	observations []GravityObservation,
	policy ObservationWindowPolicy,
) ObservationWindow {
	if policy.MaximumStableRise <= 0 {
		policy.MaximumStableRise = 0.010
	}
	return ObservationWindow{
		observations: append([]GravityObservation(nil), observations...),
		policy:       policy,
	}
}

func (w ObservationWindow) ObservationCount() int {
	return len(w.metricObservations())
}

func (w ObservationWindow) FinalObservation() (GravityObservation, bool) {
	observations := w.metricObservations()
	if len(observations) == 0 {
		return GravityObservation{}, false
	}
	return observations[len(observations)-1], true
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

func (w ObservationWindow) metricObservations() []GravityObservation {
	observations := append([]GravityObservation(nil), w.observations...)
	if !w.policy.TrimTrailingRise || len(observations) < 2 {
		return observations
	}
	maximumRise := w.policy.MaximumStableRise
	if maximumRise <= 0 {
		maximumRise = 0.010
	}
	end := len(observations)
	for end > 1 {
		previous := observations[end-2]
		latest := observations[end-1]
		rise := latest.Gravity - previous.Gravity
		if rise <= 0 || rise > maximumRise {
			break
		}
		end--
	}
	if w.policy.MinimumSamples > 0 && end < w.policy.MinimumSamples {
		return observations
	}
	return observations[:end]
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
		FinalGravity:        round2(finalGravity),
		ApparentAttenuation: round2(attenuation),
		EstimatedAlcoholABV: round2(alcohol),
		ObservationCount:    count,
	}
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}
