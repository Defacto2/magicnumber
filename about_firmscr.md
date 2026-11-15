# The FiRM `FIRM.SCR` EGA RLE Format

`FIRM.SCR` is a boutique, headerless RLE-compressed EGA graphic sidecar format used by late-1989 MS-DOS cracktros and game loaders from **The FiRM** (e.g., *Operation Wolf* and *Licence To Kill*, with loader code by *The Software Surgeon* and artwork credited to *Lord IronWolffe*).

Unlike standard DOS paint/capture formats (PCX, ColoRIX, Dr. Halo, Deluxe Paint, EGA Paint, or BSAVE), `FIRM.SCR` has no magic bytes, dimensions, or format header. No `magicnumber` signature is registered for it.

## Observed Samples

| Sample | Encoded Size | RLE EOF Clean? | Decoded Size (`rep_add=1`) | Decoded Size (`rep_add=0`) |
| :--- | :--- | :--- | :--- | :--- |
| `output.000.FIRM.SCR` | `29,070` bytes (`0x718e`) | Yes (`29070/29070`) | `35,251` bytes (`0x89b3`) | `34,143` bytes (`0x855f`) |
| `output.001.FIRM.SCR` | `18,289` bytes (`0x4771`) | Yes (`18289/18289`) | `27,500` bytes (`0x6b6c`) | `26,086` bytes (`0x65e6`) |

## RLE Stream Encoding

The file starts immediately at offset `0` with a byte-oriented high-bit RLE stream:

* **`b & 0x80 != 0` (repeat run):** Repeat the next byte `(b & 0x7F) + 1` times (`rep_add=1`), or `(b & 0x7F)` times (`rep_add=0`).
* **`b & 0x80 == 0` (literal run):** Copy the next `b` bytes verbatim (`lit_add=0`).

Both observed samples decode cleanly to exact EOF under this framing (`lit_add=0`).

In `output.000.FIRM.SCR`, the first two bytes are `80 02`, followed by `10` (16 literal bytes):

```text
80 02 10 01 00 14 03 04 05 27 07 38 39 3a 3b 3c 36 06 3f ff 00 ...
```

Under `rep_add=1`, `80 02` emits a single `0x02` byte followed by the 16-byte literal run (`01 00 14 03 04 05 27 07 38 39 3a 3b 3c 36 06 3f`). (Under `rep_add=0`, `80 02` would emit 0 bytes.)

## Decoded Prefix (EGA Palette)

Under `rep_add=1`, both decoded streams begin with a 1-byte prefix followed by a 16-byte 6-bit EGA hardware palette (values `0x00`–`0x3F`), then variable-length screen data:

* `output.000.FIRM.SCR`: `02` + `01 00 14 03 04 05 27 07 38 39 3a 3b 3c 36 06 3f`
* `output.001.FIRM.SCR`: `00` + `00 00 00 00 00 00 00 07 11 39 00 00 34 3d 00 3f`

## Loader (`output.001.OPWOLF.COM`)

The accompanying loader [`output.001.OPWOLF.COM`] (`4,065` bytes, dated `Sep 16 1989`) is an `MZ` MS-DOS `.EXE` with a `.COM` extension, packed with **Microsoft EXEPACK** (unpacking to `5,184` bytes in memory) and containing the banner:

```text
  FiRMLOADER-2   was written by       The       Software Surgeons
```

Static inspection of the unpacked entry stub shows:

* **Video probing:** Checks for VGA via `INT 10h AX=1A00h` (`BL == 7` or `8`) and EGA via `INT 10h AX=1200h BL=10h` (`CL & 0x0F == 0x09`, enhanced EGA color monitor switch setting).
* **State saving & timing:** Saves the caller's video mode (`INT 10h AH=0Fh`), cursor (`INT 10h AH=03h`), interrupt vectors, and BIOS data area, and runs a timer calibration loop.
* **Obfuscation / anti-debug:** Hooks `INT 23h` and `INT 03h` (using `INT 03h` as a trap/call mechanism) and encrypts the main loader body and strings (including the `FIRM.SCR` filename, file I/O, RLE decoder, and EGA VRAM writes). Because of that encryption, the exact authoring tool or EGA drawing routine cannot be determined from static disassembly alone without tracing/emulating the 16-bit loader.

## Why It Does Not Have a `magicnumber` Signature

`magicnumber.Find` inspects content bytes only (without filenames) and currently reports both `FIRM.SCR` samples as `Binary data or binary text`. Adding a standalone content signature is risky without strong structural invariants:

1. **No magic header:** The stream begins directly with RLE control bytes.
2. **Variable decoded size:** Unlike `SPICGADump` (which decodes to a fixed `16,384`-byte CGA frame), decoded sizes vary (`35,251` vs `27,500` bytes).
3. **False-positive risk:** Unconstrained `0x80`-flagged RLE can occasionally walk unrelated binary files to EOF. A safe detector would need to stack exact-EOF RLE decoding, size bounds, non-text checks, and 6-bit EGA palette validation (`<= 0x3F` for the first 17 decoded bytes) and be verified against a broader binary corpus.
