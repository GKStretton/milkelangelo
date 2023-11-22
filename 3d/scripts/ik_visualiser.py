import bpy
import math

YAW_ZERO_OFFSET = 26

PRECISION = 5  # Decimal point precision

class Circle(object):
    """ An OOP implementation of a circle as an object """

    def __init__(self, xposition, yposition, radius):
        self.xpos = xposition
        self.ypos = yposition
        self.radius = radius

    def circle_intersect(self, circle2):
        """
        Intersection points of two circles using the construction of triangles
        as proposed by Paul Bourke, 1997.
        http://paulbourke.net/geometry/circlesphere/
        """
        X1, Y1 = self.xpos, self.ypos
        X2, Y2 = circle2.xpos, circle2.ypos
        R1, R2 = self.radius, circle2.radius

        Dx = X2-X1
        Dy = Y2-Y1
        D = round(math.sqrt(Dx**2 + Dy**2), PRECISION)
        # Distance between circle centres
        if D > R1 + R2:
            return "The circles do not intersect"
        elif D < math.fabs(R2 - R1):
            return "No Intersect - One circle is contained within the other"
        elif D == 0 and R1 == R2:
            return "No Intersect - The circles are equal and coincident"
        else:
            if D == R1 + R2 or D == R1 - R2:
                CASE = "The circles intersect at a single point"
            else:
                CASE = "The circles intersect at two points"
            chorddistance = (R1**2 - R2**2 + D**2)/(2*D)
            # distance from 1st circle's centre to the chord between intersects
            halfchordlength = math.sqrt(R1**2 - chorddistance**2)
            chordmidpointx = X1 + (chorddistance*Dx)/D
            chordmidpointy = Y1 + (chorddistance*Dy)/D
            I1 = (round(chordmidpointx + (halfchordlength*Dy)/D, PRECISION),
                  round(chordmidpointy - (halfchordlength*Dx)/D, PRECISION))
            theta1 = round(math.degrees(math.atan2(I1[1]-Y1, I1[0]-X1)),
                           PRECISION)
            I2 = (round(chordmidpointx - (halfchordlength*Dy)/D, PRECISION),
                  round(chordmidpointy + (halfchordlength*Dx)/D, PRECISION))
            theta2 = round(math.degrees(math.atan2(I2[1]-Y1, I2[0]-X1)),
                           PRECISION)
            if theta2 > theta1:
                I1, I2 = I2, I1
            return (I1, I2, CASE)


debugFormat = """
DATA
x: {:.2f} y: {:.2f}
ring1: {:.2f}deg
ring2: {:.2f}deg
ringChoice: {:.2f}deg
yawOffset: {:.2f}deg
"""

def distanceBetweenPoints(x1, y1, x2, y2):
    return math.sqrt((x1 - x2)**2+(y1-y2)**2)

class ModalTimerOperator(bpy.types.Operator):
    """Operator which runs itself from a timer"""
    bl_idname = "wm.modal_timer_operator"
    bl_label = "Modal Timer Operator"

    _timer = None
    previousRing = 0

    def modal(self, context, event):
        if event.type in {'RIGHTMOUSE', 'ESC'}:
            self.cancel(context)
            return {'CANCELLED'}

        if event.type == 'TIMER':
            stage_radius = 75
            x_raw = bpy.data.objects["TipTarget"].location[0]
            y_raw = bpy.data.objects["TipTarget"].location[1]
            x = x_raw / stage_radius
            y = y_raw / stage_radius

            armAxisEmpty = bpy.data.objects["Axis Centre"]
            armCentre = (armAxisEmpty.matrix_world.translation[0], armAxisEmpty.matrix_world.translation[1])
            armPathRadius = distanceBetweenPoints(0, 0, armCentre[0], armCentre[1])

            print("armPathRadius", armPathRadius)
            
            C1 = Circle(0, 0, armPathRadius)
            C2 = Circle(x_raw, y_raw, armPathRadius)

            (i1, i2, case) = C1.circle_intersect(C2)

            print("i1 distance from centre", distanceBetweenPoints(0,0,i1[0],i1[1]))

            bpy.data.objects["RingIntersect1"].location[0] = i1[0]
            bpy.data.objects["RingIntersect1"].location[1] = i1[1]
            bpy.data.objects["RingIntersect2"].location[0] = i2[0]
            bpy.data.objects["RingIntersect2"].location[1] = i2[1]

            # Now we have the 2 intersect points, and target x_raw, y_raw

            angle1 = (-math.atan2(i1[0], i1[1]) * 180 / math.pi + 270) % 360
            angle2 = (-math.atan2(i2[0], i2[1]) * 180 / math.pi + 270) % 360

            def ringValid(a):
                if a >= 0 and a <= 215:
                    return True
                return False

            intersection = i1
            if ringValid(angle1) and ringValid(angle2):
                # move to whichever is closest to previous position
                if abs(angle1 - self.previousRing) < abs(angle2 - self.previousRing):
                    ring = angle1
                else:
                    intersection = i2
                    ring = angle2
            elif ringValid(angle1):
                ring = angle1
            elif ringValid(angle2):
                intersection = i2
                ring = angle2
            else:
                ring = -13
            
            def angleBetweenVectors(x1, y1, x2, y2):
                print("v1:", x1, y1)
                print("v2:", x2, y2)
                dot = x1*x2 + y1*y2      # dot product between [x1, y1] and [x2, y2]
                det = x1*y2 - y1*x2      # determinant
                angle = math.atan2(det, dot) * 180 / math.pi
                print("angle:", angle)
                return angle

            yawOffset = -angleBetweenVectors(-intersection[0], -intersection[1], x_raw - intersection[0], y_raw - intersection[1])
            yaw = yawOffset + YAW_ZERO_OFFSET
            
            bpy.data.objects["DebugText"].data.body = debugFormat.format(x, y, angle1, angle2, ring, yawOffset)
            
            bpy.data.scenes["Scene"].ring = ring
            bpy.data.scenes["Scene"].pipette_yaw = yaw

            self.previousRing = ring

        return {'PASS_THROUGH'}

    def execute(self, context):
        wm = context.window_manager
        self._timer = wm.event_timer_add(0.1, window=context.window)
        wm.modal_handler_add(self)
        return {'RUNNING_MODAL'}

    def cancel(self, context):
        wm = context.window_manager
        wm.event_timer_remove(self._timer)


def menu_func(self, context):
    self.layout.operator(ModalTimerOperator.bl_idname, text=ModalTimerOperator.bl_label)


def register():
    bpy.utils.register_class(ModalTimerOperator)
    bpy.types.VIEW3D_MT_view.append(menu_func)


# Register and add to the "view" menu (required to also use F3 search "Modal Timer Operator" for quick access).
def unregister():
    bpy.utils.unregister_class(ModalTimerOperator)
    bpy.types.VIEW3D_MT_view.remove(menu_func)


if __name__ == "__main__":
    register()

    # test call
    bpy.ops.wm.modal_timer_operator()
