def bytes_to_crumbs(bytes):
    for byte in bytes:
        for _ in range(4):
            yield byte & 0b11
            byte >>= 2


def debiase_crumbs(crumbs):
    for crumb in crumbs:
        if crumb == 0b10:
            yield 1
        elif crumb == 0b01:
            yield 0


def bits_to_bytes(bits):
    """Just discard."""
    try:
        while True:
            byte = next(bits)
            for _ in range(7):
                byte <<= 1
                byte |= next(bits)
            yield byte
    except StopIteration:
        pass


def bits_to_bytes(bits):
    """Pad to the right - 1111000."""
    try:
        while True:
            byte = next(bits)
            for _ in range(7):
                byte <<= 1
                byte |= next(bits, 0)
            yield byte
    except StopIteration:
        pass


def bits_to_bytes(bits):
    """Pad to the left - 00001111."""
    try:
        while True:
            byte = next(bits)
            for _ in range(7):
                try:
                    bit = next(bits)
                except StopIteration:
                    yield byte
                    return
                byte <<= 1
                byte |= bit
    except StopIteration:
        pass


def get_debiased_bytes(length, path):
    with open(path, 'rb') as f:
        _bytes = f.read(length)
    return bytes(bits_to_bytes(debiase_crumbs(bytes_to_crumbs(_bytes))))

directory = "samples"

import os

for filename in os.listdir(directory):
    if filename.endswith(".wav"):
        path = os.path.join(directory, filename)
        print(path, len(get_debiased_bytes(-1, path)))