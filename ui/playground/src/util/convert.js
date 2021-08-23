/**
 * Convert a freeNorm in one range into a freeNorm in another. Performs a linear transformation.
 *
 * @param {number} oldMin
 * @param {number} oldMax
 * @param {number} newMin
 * @param {number} newMax
 * @param {number} value
 * @return {number}
 */
export function linearConversion(oldMin, oldMax, newMin, newMax, value) {
  const oldSize = oldMax - oldMin;
  const newSize = newMax - newMin;
  return newMin + (value - oldMin) * newSize / oldSize;
}
