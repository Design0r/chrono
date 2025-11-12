export function hexToHSL(hex: string): [number, number, number] {
  // Remove leading "#" if present
  const clean = hex.startsWith("#") ? hex.slice(1) : hex;

  // Must be exactly 6 hex chars
  if (clean.length !== 6) return [0, 0, 0];

  const rgbValue = parseInt(clean, 16);
  if (Number.isNaN(rgbValue)) return [0, 0, 0];

  // Extract 0..255 channels, normalize to 0..1
  const r = ((rgbValue >> 16) & 0xff) / 255;
  const g = ((rgbValue >> 8) & 0xff) / 255;
  const b = (rgbValue & 0xff) / 255;

  const maxVal = Math.max(r, g, b);
  const minVal = Math.min(r, g, b);
  const l = (maxVal + minVal) / 2;

  let h = 0;
  let s = 0;

  if (maxVal !== minVal) {
    const d = maxVal - minVal;

    // Saturation
    s = l > 0.5 ? d / (2 - maxVal - minVal) : d / (maxVal + minVal);

    // Hue
    switch (maxVal) {
      case r:
        h = (g - b) / d + (g < b ? 6 : 0);
        break;
      case g:
        h = (b - r) / d + 2;
        break;
      case b:
        h = (r - g) / d + 4;
        break;
    }
    h /= 6;
  }

  const hDegrees = h * 360;
  return [hDegrees, s, l];
}

export function hsla(h: number, s: number, l: number, a: number): string {
  return `hsla(${h.toFixed(1)}, ${(s * 100).toFixed(1)}%, ${l * 100}%, ${a})`;
}
