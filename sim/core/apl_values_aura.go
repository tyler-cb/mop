package core

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
)

type APLValueAuraIsKnown struct {
	DefaultAPLValueImpl
	aura AuraReference
}

func (rot *APLRotation) newValueAuraIsKnown(config *proto.APLValueAuraIsKnown, _ *proto.UUID) APLValue {
	aura := rot.GetAPLAura(rot.GetSourceUnit(config.SourceUnit), config.AuraId)
	return &APLValueAuraIsKnown{
		aura: aura,
	}
}
func (value *APLValueAuraIsKnown) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAuraIsKnown) GetBool(sim *Simulation) bool {
	return value.aura.Get() != nil
}
func (value *APLValueAuraIsKnown) String() string {
	return fmt.Sprintf("Aura Active(%s)", value.aura.String())
}

type APLValueAuraIsActive struct {
	DefaultAPLValueImpl
	aura AuraReference
}

func (rot *APLRotation) newValueAuraIsActive(config *proto.APLValueAuraIsActive, _ *proto.UUID) APLValue {
	if config.AuraId == nil {
		return nil
	}
	aura := rot.GetAPLAura(rot.GetSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLValueAuraIsActive{
		aura: aura,
	}
}
func (value *APLValueAuraIsActive) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAuraIsActive) GetBool(sim *Simulation) bool {
	return value.aura.Get().IsActive()
}
func (value *APLValueAuraIsActive) String() string {
	return fmt.Sprintf("Aura Active(%s)", value.aura.String())
}

type APLValueAuraIsActiveWithReactionTime struct {
	DefaultAPLValueImpl
	aura         AuraReference
	reactionTime time.Duration
}

func (rot *APLRotation) newValueAuraIsActiveWithReactionTime(config *proto.APLValueAuraIsActiveWithReactionTime, _ *proto.UUID) APLValue {
	if config.AuraId == nil {
		return nil
	}
	aura := rot.GetAPLAura(rot.GetSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLValueAuraIsActiveWithReactionTime{
		aura:         aura,
		reactionTime: rot.unit.ReactionTime,
	}
}
func (value *APLValueAuraIsActiveWithReactionTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAuraIsActiveWithReactionTime) GetBool(sim *Simulation) bool {
	aura := value.aura.Get()
	return aura.IsActive() && aura.TimeActive(sim) >= value.reactionTime
}
func (value *APLValueAuraIsActiveWithReactionTime) String() string {
	return fmt.Sprintf("Aura Active With Reaction Time(%s)", value.aura.String())
}

type APLValueAuraIsInactiveWithReactionTime struct {
	DefaultAPLValueImpl
	aura         AuraReference
	reactionTime time.Duration
}

func (rot *APLRotation) newValueAuraIsInactiveWithReactionTime(config *proto.APLValueAuraIsInactiveWithReactionTime, _ *proto.UUID) APLValue {
	if config.AuraId == nil {
		return nil
	}
	aura := rot.GetAPLAura(rot.GetSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLValueAuraIsInactiveWithReactionTime{
		aura:         aura,
		reactionTime: rot.unit.ReactionTime,
	}
}
func (value *APLValueAuraIsInactiveWithReactionTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAuraIsInactiveWithReactionTime) GetBool(sim *Simulation) bool {
	aura := value.aura.Get()
	return !aura.IsActive() && aura.TimeInactive(sim) >= value.reactionTime
}
func (value *APLValueAuraIsInactiveWithReactionTime) String() string {
	return fmt.Sprintf("Aura Inactive With Reaction Time(%s)", value.aura.String())
}

type APLValueAuraRemainingTime struct {
	DefaultAPLValueImpl
	aura AuraReference
}

func (rot *APLRotation) newValueAuraRemainingTime(config *proto.APLValueAuraRemainingTime, _ *proto.UUID) APLValue {
	if config.AuraId == nil {
		return nil
	}
	aura := rot.GetAPLAura(rot.GetSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLValueAuraRemainingTime{
		aura: aura,
	}
}
func (value *APLValueAuraRemainingTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueAuraRemainingTime) GetDuration(sim *Simulation) time.Duration {
	aura := value.aura.Get()
	return TernaryDuration(aura.IsActive(), aura.RemainingDuration(sim), 0)
}
func (value *APLValueAuraRemainingTime) String() string {
	return fmt.Sprintf("Aura Remaining Time(%s)", value.aura.String())
}

type APLValueAuraNumStacks struct {
	DefaultAPLValueImpl
	aura AuraReference
}

func (rot *APLRotation) newValueAuraNumStacks(config *proto.APLValueAuraNumStacks, uuid *proto.UUID) APLValue {
	if config.AuraId == nil {
		return nil
	}
	aura := rot.GetAPLAura(rot.GetSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	if aura.Get().MaxStacks == 0 {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Warning, "%s is not a stackable aura", ProtoToActionID(config.AuraId))
		return nil
	}
	return &APLValueAuraNumStacks{
		aura: aura,
	}
}
func (value *APLValueAuraNumStacks) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueAuraNumStacks) GetInt(sim *Simulation) int32 {
	return value.aura.Get().GetStacks()
}
func (value *APLValueAuraNumStacks) String() string {
	return fmt.Sprintf("Aura Num Stacks(%s)", value.aura.String())
}

type APLValueAuraInternalCooldown struct {
	DefaultAPLValueImpl
	aura AuraReference
}

func (rot *APLRotation) newValueAuraInternalCooldown(config *proto.APLValueAuraInternalCooldown, _ *proto.UUID) APLValue {
	if config.AuraId == nil {
		return nil
	}
	aura := rot.GetAPLICDAura(rot.GetSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLValueAuraInternalCooldown{
		aura: aura,
	}
}
func (value *APLValueAuraInternalCooldown) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueAuraInternalCooldown) GetDuration(sim *Simulation) time.Duration {
	return value.aura.Get().Icd.TimeToReady(sim)
}
func (value *APLValueAuraInternalCooldown) String() string {
	return fmt.Sprintf("Aura Remaining ICD(%s)", value.aura.String())
}

type APLValueAuraICDIsReadyWithReactionTime struct {
	DefaultAPLValueImpl
	aura         AuraReference
	reactionTime time.Duration
}

func (rot *APLRotation) newValueAuraICDIsReadyWithReactionTime(config *proto.APLValueAuraICDIsReadyWithReactionTime, _ *proto.UUID) APLValue {
	if config.AuraId == nil {
		return nil
	}
	aura := rot.GetAPLICDAura(rot.GetSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLValueAuraICDIsReadyWithReactionTime{
		aura:         aura,
		reactionTime: rot.unit.ReactionTime,
	}
}
func (value *APLValueAuraICDIsReadyWithReactionTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAuraICDIsReadyWithReactionTime) GetBool(sim *Simulation) bool {
	aura := value.aura.Get()
	return aura.Icd.IsReady(sim) || (aura.IsActive() && aura.TimeActive(sim) < value.reactionTime)
}
func (value *APLValueAuraICDIsReadyWithReactionTime) String() string {
	return fmt.Sprintf("Aura ICD Is Ready with Reaction Time(%s)", value.aura.String())
}

type APLValueAuraShouldRefresh struct {
	DefaultAPLValueImpl
	aura       AuraReference
	maxOverlap APLValue
}

func (rot *APLRotation) newValueAuraShouldRefresh(config *proto.APLValueAuraShouldRefresh, uuid *proto.UUID) APLValue {
	if config.AuraId == nil {
		return nil
	}
	aura := rot.GetAPLAura(rot.GetTargetUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}

	maxOverlap := rot.coerceTo(rot.newAPLValue(config.MaxOverlap), proto.APLValueType_ValueTypeDuration)
	if maxOverlap == nil {
		maxOverlap = rot.newValueConst(&proto.APLValueConst{Val: "0ms"}, uuid)
	}

	return &APLValueAuraShouldRefresh{
		aura:       aura,
		maxOverlap: maxOverlap,
	}
}
func (value *APLValueAuraShouldRefresh) GetInnerValues() []APLValue {
	return []APLValue{value.maxOverlap}
}
func (value *APLValueAuraShouldRefresh) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAuraShouldRefresh) GetBool(sim *Simulation) bool {
	return value.aura.Get().ShouldRefreshExclusiveEffects(sim, value.maxOverlap.GetDuration(sim))
}
func (value *APLValueAuraShouldRefresh) String() string {
	return fmt.Sprintf("Should Refresh Aura(%s)", value.aura.String())
}
