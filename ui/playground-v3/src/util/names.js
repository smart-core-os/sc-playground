/**
 * @param {string} name
 * @return {string}
 */
export function localName(name) {
  if (!name) return name;
  const lastDot = name.lastIndexOf('.');
  return name.substring(lastDot + 1)
}
