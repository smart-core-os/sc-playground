/**
 * @param {ElectricMode.Segment.AsObject[]} segments
 * @return {number}
 */

export function maxMagnitude(segments) {
  return segments.reduce((max, s) => Math.max(max, s.magnitude), 0);
}
