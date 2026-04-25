package candy_physics

import "candy/candy_evaluator"

// RegisterBuiltins registers all physics builtins with the Candy evaluator.
func RegisterBuiltins() {
	// World
	candy_evaluator.RegisterBuiltin("physicsCreateWorld", builtinPhysicsCreateWorld)
	candy_evaluator.RegisterBuiltin("physicsDestroyWorld", builtinPhysicsDestroyWorld)
	candy_evaluator.RegisterBuiltin("physicsStep", builtinPhysicsStep)
	candy_evaluator.RegisterBuiltin("physicsSetGravity", builtinPhysicsSetGravity)
	candy_evaluator.RegisterBuiltin("physicsGetGravity", builtinPhysicsGetGravity)

	// Body creation
	candy_evaluator.RegisterBuiltin("physicsCreateBox", builtinPhysicsCreateBox)
	candy_evaluator.RegisterBuiltin("physicsCreateSphere", builtinPhysicsCreateSphere)
	candy_evaluator.RegisterBuiltin("physicsCreateCapsule", builtinPhysicsCreateCapsule)
	candy_evaluator.RegisterBuiltin("physicsCreatePlane", builtinPhysicsCreatePlane)
	candy_evaluator.RegisterBuiltin("physicsDestroyBody", builtinPhysicsDestroyBody)

	// Position / velocity / rotation
	candy_evaluator.RegisterBuiltin("physicsGetPosition", builtinPhysicsGetPosition)
	candy_evaluator.RegisterBuiltin("physicsSetPosition", builtinPhysicsSetPosition)
	candy_evaluator.RegisterBuiltin("physicsGetVelocity", builtinPhysicsGetVelocity)
	candy_evaluator.RegisterBuiltin("physicsSetVelocity", builtinPhysicsSetVelocity)
	candy_evaluator.RegisterBuiltin("physicsGetAngularVelocity", builtinPhysicsGetAngularVelocity)
	candy_evaluator.RegisterBuiltin("physicsSetAngularVelocity", builtinPhysicsSetAngularVelocity)
	candy_evaluator.RegisterBuiltin("physicsGetRotation", builtinPhysicsGetRotation)
	candy_evaluator.RegisterBuiltin("physicsSetRotation", builtinPhysicsSetRotation)
	candy_evaluator.RegisterBuiltin("physicsSetRotationFromAxisAngle", builtinPhysicsSetRotationFromAxisAngle)

	// Forces / impulses
	candy_evaluator.RegisterBuiltin("physicsApplyForce", builtinPhysicsApplyForce)
	candy_evaluator.RegisterBuiltin("physicsApplyForceAtPoint", builtinPhysicsApplyForceAtPoint)
	candy_evaluator.RegisterBuiltin("physicsApplyImpulse", builtinPhysicsApplyImpulse)
	candy_evaluator.RegisterBuiltin("physicsApplyImpulseAtPoint", builtinPhysicsApplyImpulseAtPoint)
	candy_evaluator.RegisterBuiltin("physicsApplyTorque", builtinPhysicsApplyTorque)

	// Material properties
	candy_evaluator.RegisterBuiltin("physicsSetMass", builtinPhysicsSetMass)
	candy_evaluator.RegisterBuiltin("physicsGetMass", builtinPhysicsGetMass)
	candy_evaluator.RegisterBuiltin("physicsSetRestitution", builtinPhysicsSetRestitution)
	candy_evaluator.RegisterBuiltin("physicsSetFriction", builtinPhysicsSetFriction)
	candy_evaluator.RegisterBuiltin("physicsSetLinearDrag", builtinPhysicsSetLinearDrag)
	candy_evaluator.RegisterBuiltin("physicsSetAngularDrag", builtinPhysicsSetAngularDrag)

	// State
	candy_evaluator.RegisterBuiltin("physicsSetActive", builtinPhysicsSetActive)
	candy_evaluator.RegisterBuiltin("physicsIsActive", builtinPhysicsIsActive)
	candy_evaluator.RegisterBuiltin("physicsIsSleeping", builtinPhysicsIsSleeping)
	candy_evaluator.RegisterBuiltin("physicsWakeBody", builtinPhysicsWakeBody)
	candy_evaluator.RegisterBuiltin("physicsSetUserData", builtinPhysicsSetUserData)
	candy_evaluator.RegisterBuiltin("physicsGetUserData", builtinPhysicsGetUserData)

	// Queries
	candy_evaluator.RegisterBuiltin("physicsGetBodyCount", builtinPhysicsGetBodyCount)
	candy_evaluator.RegisterBuiltin("physicsGetContacts", builtinPhysicsGetContacts)

	// Raycasting
	candy_evaluator.RegisterBuiltin("physicsCastRay", builtinPhysicsCastRay)
	candy_evaluator.RegisterBuiltin("physicsCastRayFirst", builtinPhysicsCastRayFirst)

	// Motion type constants (call with () to get the int value)
	candy_evaluator.RegisterBuiltin("PHYSICS_STATIC", builtinPhysicsMotionStatic)
	candy_evaluator.RegisterBuiltin("PHYSICS_DYNAMIC", builtinPhysicsMotionDynamic)
	candy_evaluator.RegisterBuiltin("PHYSICS_KINEMATIC", builtinPhysicsMotionKinematic)
}
