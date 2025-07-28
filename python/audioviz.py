from matplotlib import pyplot as plt
import pathlib

audio_bin = pathlib.Path("audio.bin").read_bytes()
plt.plot(audio_bin)
plt.show()