/**
 * @param {ElectricMode.Segment.AsObject[]} segments
 * @return {number}
 */

export function maxMagnitude(segments) {
  return segments.reduce((max, s) => Math.max(max, s.magnitude), 0);
}


/**
 * Convert a timestamp object to a Date. Timestamp.AsObject doesn't have this method :(
 *
 * @param {Timestamp.AsObject} ts
 * @return {Date}
 */
export function toDate(ts) {
  return new Date((ts.seconds * 1000) + (ts.nanos / 1000000));
}

/**
 * Convert a duration object into a millisecond value.
 *
 * @param {Duration.AsObject} d
 * @return {number}
 */
export function durationMillis(d) {
  return d.seconds * 1000 + d.nanos / 1_000_000;
}

/**
 * Convert a duration into a string of the format 12H3M14S etc.
 * @param {number|Duration.AsObject} d Either numerical milliseconds or a Duration.AsObject.
 */
export function durationString(d) {
  if (typeof d === 'object') {
    d = durationMillis(d)
  }
  const millis = 1;
  const seconds = 1000 * millis;
  const minutes = 60 * seconds;
  const hours = 60 * minutes;
  const days = 24 * hours;
  const units = [
    {n: 'D', m: days},
    {n: 'H', m: hours},
    {n: 'M', m: minutes},
    {n: 'S', m: seconds},
    {n: 'MS', m: millis},
  ]

  const parts = [];
  for (const unit of units) {
    if (d > unit.m) {
      const v = Math.floor(d / unit.m)
      d -= v * unit.m
      parts.push(`${v}${unit.n}`)
    }
  }

  return parts.join('')
}
