#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

uniform sampler2D texture0;
uniform vec4 colDiffuse;

out vec4 finalColor;

uniform float darkenStrength = 0.5;

void main()
{
    vec4 texColor = texture(texture0, fragTexCoord);

    vec3 darkened = texColor.rgb * darkenStrength;

    finalColor = vec4(darkened, texColor.a) * fragColor * colDiffuse;
}