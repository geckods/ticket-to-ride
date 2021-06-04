import json
from PIL import Image, ImageDraw, ImageFont
import graphviz

with open("game.log", 'r') as f:
    lines = f.readlines()

currText = ""
currGraph = ""
currGraphLoc = "basegraph.png"

image = Image.new('RGB', (2439, 1600), (255, 255, 255))

frames = []
i = 0
for line in lines:
    print("Photo", i)
    i+=1
    thelog = json.loads(line)
    if 'EVENT' in thelog and thelog['EVENT'] == "GRAPH":
        if currGraph != thelog["GRAPH"]:
            currGraph = thelog["GRAPH"]
            with open("graph.txt", "w") as f:
                f.write(currGraph)
            currGraphLoc = graphviz.render('neato', 'png', 'graph.txt')
    else:
        currText = thelog["msg"]

    newimage = Image.new('RGB', (2439, 1600), (255, 255, 255))
    graphImage = Image.open(currGraphLoc)
    newimage.paste(graphImage, (0, 0))
    drawable = ImageDraw.Draw(newimage)

    font = ImageFont.truetype("UberMoveTextRegular.otf", 30)
    drawable.text((100, 1500), currText, (0, 0, 0), font)

    # newimage.show()
    frames.append(newimage.reduce(2))
    # newimage.save("currimage.png")
    # input()

print("Done")
frames[0].save('png_to_gif2.gif', format='GIF', append_images=frames[1:], save_all=True, duration=300, loop=0, optimize=True)
