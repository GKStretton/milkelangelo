import bpy
from math import radians

bl_info = {
    # required
    'name': 'Light Tip Path Visualiser',
    'blender': (2, 93, 0),
    'category': 'Object',
}

def updateRing(self, context):
    bpy.data.objects["Pipette Movement Axis"].rotation_euler[2] = radians(self.ring)
    
def updateYaw(self, context):
    bpy.data.objects["PULLEY_80_V10"].rotation_euler[2] = -radians(self.pipette_yaw - 40)
    
def updatePitch(self, context):
    bpy.data.objects["PULLEY_60_V10"].rotation_euler[0] = radians(self.pipette_pitch)
    
def updateZ(self, context):
    bpy.data.objects["Z base"].location[2] = self.z_axis
    
def updateXY(self, context):
    bpy.data.objects["TipTarget"].location[0] = self.x_coord * 75
    bpy.data.objects["TipTarget"].location[1] = self.y_coord * 75

PROPS = [
    ('ring', bpy.props.FloatProperty(
        name='Ring Rotate',
        update=updateRing,
        soft_min=0,
        soft_max=215,
        precision=1,
        step=25,
    )),
    ('pipette_yaw', bpy.props.FloatProperty(
        name='Pipette Yaw',
        update=updateYaw,
        soft_min=-40,
        soft_max=220,
        precision=1,
        step=25,
    )),
    ('pipette_pitch', bpy.props.FloatProperty(
        name='Pipette Pitch',
        update=updatePitch,
        soft_min=0,
        soft_max=90,
        precision=1,
        step=25,
    )),
    ('z_axis', bpy.props.FloatProperty(
        name='Platform Z',
        update=updateZ,
        soft_min=0,
        soft_max=75,
        precision=1,
        step=25,
    )),
    ('x_coord', bpy.props.FloatProperty(
        name='X',
        update=updateXY,
        soft_min=-1,
        soft_max=1,
        precision=2,
        step=0.05,
    )),
    ('y_coord', bpy.props.FloatProperty(
        name='Y',
        update=updateXY,
        soft_min=-1,
        soft_max=1,
        precision=2,
        step=0.05,
    )),
]

class TestPanel(bpy.types.Panel):
    bl_idname = 'VIEW3D_PT_test_panel'
    bl_label = 'Light Visualisation'
    bl_space_type = 'VIEW_3D'
    bl_region_type = 'UI'
    
    def draw(self, context):
        col = self.layout.column()
        for (prop_name, _) in PROPS:
            row = col.row()
            row.prop(context.scene, prop_name)
        col.operator('object.reset_light', text='Reset')

class ResetLight(bpy.types.Operator):
    bl_idname = 'object.reset_light'
    bl_label = 'Reset Light'
    bl_options = {'REGISTER', 'UNDO'}
    
    def execute(self, context):
        scene = bpy.data.scenes["Scene"]
        
        scene.ring = 0
        scene.pipette_pitch = 0
        scene.pipette_yaw = 0
        scene.z_axis = 0
        
        return {'FINISHED'}
        
CLASSES = [
    TestPanel,
    ResetLight,
]

def my_handler(a, b):
    updateRing(bpy.data.scenes['Scene'], 0)
    updateYaw(bpy.data.scenes['Scene'], 0)
    updatePitch(bpy.data.scenes['Scene'], 0)

def register():
    print('registered') # just for debug
    
    for (prop_name, prop_value) in PROPS:
        setattr(bpy.types.Scene, prop_name, prop_value)
    
    for klass in CLASSES:
        bpy.utils.register_class(klass)
    
    
    bpy.app.handlers.frame_change_pre.clear()
    bpy.app.handlers.frame_change_pre.append(my_handler)

def unregister():
    print('unregistered') # just for debug
    
    for (prop_name, _) in PROPS:
        delattr(bpy.types.Scene, prop_name)
    
    for klass in CLASSES:
        bpy.utils.unregister_class(klass)
    
    bpy.app.handlers.frame_change_pre.clear()
    

        
if __name__ == "__main__":
    register()