#pragma once

#define PI_F 3.141592653589f

/* CircleCircleIntersection() *
 * Determine the points where 2 circles in a common plane intersect.
 *
 * int CircleCircleIntersection(
 *                                // center and radius of 1st circle
 *                                float x0, float y0, float r0,
 *                                // center and radius of 2nd circle
 *                                float x1, float y1, float r1,
 *                                // 1st intersection point
 *                                float *xi, float *yi,              
 *                                // 2nd intersection point
 *                                float *xi_prime, float *yi_prime)
 *
 * This is a public domain work. 3/26/2005 Tim Voght
 *
 */
int CircleCircleIntersection(float x0, float y0, float r0,
							   float x1, float y1, float r1,
							   float *xi, float *yi,
							   float *xi_prime, float *yi_prime);

// AngleBetweenVectors returns angle in degrees between two specified vectors
float AngleBetweenVectors(float x0, float y0, float x1, float y1);

void boundXYToCircle(float *x, float *y, float radius);
void boundToSignedMaximum(float *n, float range);

bool numInRange(float num, float min, float max);