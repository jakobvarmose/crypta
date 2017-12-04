function numberToString(num) {
  return num.toLocaleString('en-US', {
    maximumSignificantDigits: 3,
    minimumSignificantDigits: 3,
  });
}

function fromSize(size) {
  if (Number(size) === 1) {
    return '1 byte';
  }
  if (size < 1e3) {
    return `${size} bytes`;
  }
  if (size < 1e6) {
    return `${numberToString(size / 1e3)} kB`;
  }
  if (size < 1e9) {
    return `${numberToString(size / 1e6)} MB`;
  }
  if (size < 1e12) {
    return `${numberToString(size / 1e9)} GB`;
  }
  return `${numberToString(size / 1e12)} TB`;
}

function stringCompare(a, b) {
  return String(a).localeCompare(b, 'en-US-u-kn-true');
}

export default {
  fromSize,
  stringCompare,
};
