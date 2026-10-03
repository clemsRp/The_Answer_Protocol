from PIL import Image
import sys


def pixeliser(input_path, output_path, taille_pixel=4):
    # Ouvre l'image
    image = Image.open(input_path).convert("RGBA")

    # Réduit l'image
    largeur = image.width // taille_pixel
    hauteur = image.height // taille_pixel

    petite = image.resize(
        (largeur, hauteur),
        Image.Resampling.NEAREST
    )

    # Réagrandit sans interpolation
    pixelisee = petite.resize(
        image.size,
        Image.Resampling.NEAREST
    )

    pixelisee.save(output_path)


if __name__ == "__main__":

    pixeliser(sys.argv[1], sys.argv[2], int(sys.argv[3]))
