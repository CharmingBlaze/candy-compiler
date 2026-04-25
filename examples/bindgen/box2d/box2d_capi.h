// Minimal Box2D-style C API fixture for candywrap demos.
typedef struct b2World* b2WorldId;
typedef struct b2Body* b2BodyId;

b2WorldId b2World_Create(float gx, float gy);
void b2World_Destroy(b2WorldId world);
b2BodyId b2World_CreateBody(b2WorldId world, float x, float y);
void b2Body_Destroy(b2BodyId body);
