#version 330

// Attributs d'entrée passés par raylib
in vec2 fragTexCoord;
in vec4 fragColor;

// Uniforms obligatoires de raylib
uniform sampler2D texture0;
uniform vec4 colDiffuse;

// Propriétés du flou (taille de l'écran ou de la texture)
uniform float renderWidth = 800.0;
uniform float renderHeight = 450.0;

out vec4 finalColor;

uniform float blurStrength = 3.0;

void main()
{
    vec2 texelSize = (1.0 / vec2(renderWidth, renderHeight)) * blurStrength;
    vec4 sum = vec4(0.0);

    // Noyau de flou simple (Box blur / Flou Gaussien simplifié 3x3)
    sum += texture(texture0, fragTexCoord + vec2(-1.0, -1.0) * texelSize) * 0.0625;
    sum += texture(texture0, fragTexCoord + vec2( 0.0, -1.0) * texelSize) * 0.125;
    sum += texture(texture0, fragTexCoord + vec2( 1.0, -1.0) * texelSize) * 0.0625;
    
    sum += texture(texture0, fragTexCoord + vec2(-1.0,  0.0) * texelSize) * 0.125;
    sum += texture(texture0, fragTexCoord + vec2( 0.0,  0.0) * texelSize) * 0.25;
    sum += texture(texture0, fragTexCoord + vec2( 1.0,  0.0) * texelSize) * 0.125;
    
    sum += texture(texture0, fragTexCoord + vec2(-1.0,  1.0) * texelSize) * 0.0625;
    sum += texture(texture0, fragTexCoord + vec2( 0.0,  1.0) * texelSize) * 0.125;
    sum += texture(texture0, fragTexCoord + vec2( 1.0,  1.0) * texelSize) * 0.0625;

    finalColor = sum * fragColor * colDiffuse;
}