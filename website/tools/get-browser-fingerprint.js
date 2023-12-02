const getBrowserFingerprint = ({
  hardwareOnly = false,
  enableWebgl = false,
  enableScreen = true,
  debug = false,
} = {}) => {
  const {
    cookieEnabled,
    deviceMemory,
    doNotTrack,
    hardwareConcurrency,
    language,
    languages,
    maxTouchPoints,
    platform,
    userAgent,
    vendor,
  } = window.navigator;

  const { width, height, colorDepth, pixelDepth } = enableScreen ? window.screen : {}; // undefined will remove this from the stringify down here
  const timezoneOffset = new Date().getTimezoneOffset();
  const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  const touchSupport = 'ontouchstart' in window;
  const devicePixelRatio = window.devicePixelRatio;

  const canvas = getCanvasID(debug);
  const webgl = enableWebgl ? getWebglID(debug) : undefined; // undefined will remove this from the stringify down here
  const webglInfo = enableWebgl ? getWebglInfo(debug) : undefined; // undefined will remove this from the stringify down here

  const data = hardwareOnly
    ? JSON.stringify({
      canvas,
      colorDepth,
      deviceMemory,
      devicePixelRatio,
      hardwareConcurrency,
      height,
      maxTouchPoints,
      pixelDepth,
      platform,
      touchSupport,
      webgl,
      webglInfo,
      width,
    })
    : JSON.stringify({
      canvas,
      colorDepth,
      cookieEnabled,
      deviceMemory,
      devicePixelRatio,
      doNotTrack,
      hardwareConcurrency,
      height,
      language,
      languages,
      maxTouchPoints,
      pixelDepth,
      platform,
      timezone,
      timezoneOffset,
      touchSupport,
      userAgent,
      vendor,
      webgl,
      webglInfo,
      width,
    });

  const datastring = JSON.stringify(data, null, 4);

  if (debug) console.log('fingerprint data', datastring);

  const result = sha256(datastring);
  return result;
};

export const getCanvasID = (debug) => {
  try {
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');
    const text = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`~1!2@3#4$5%6^7&8*9(0)-_=+[{]}|;:',<.>/?";
    ctx.textBaseline = 'top';
    ctx.font = "14px 'Arial'";
    ctx.textBaseline = 'alphabetic';
    ctx.fillStyle = '#f60';
    ctx.fillRect(125, 1, 62, 20);
    ctx.fillStyle = '#069';
    ctx.fillText(text, 2, 15);
    ctx.fillStyle = 'rgba(102, 204, 0, 0.7)';
    ctx.fillText(text, 4, 17);

    const result = canvas.toDataURL();

    if (debug) {
      document.body.appendChild(canvas);
    } else {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
    }

    return murmurhash3_32_gc(result);
  } catch {
    return null;
  }
};

export const getWebglID = (debug) => {
  try {
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('webgl');
    canvas.width = 256;
    canvas.height = 128;

    const f =
      'attribute vec2 attrVertex;varying vec2 varyinTexCoordinate;uniform vec2 uniformOffset;void main(){varyinTexCoordinate=attrVertex+uniformOffset;gl_Position=vec4(attrVertex,0,1);}';
    const g =
      'precision mediump float;varying vec2 varyinTexCoordinate;void main() {gl_FragColor=vec4(varyinTexCoordinate,0,1);}';
    const h = ctx.createBuffer();

    ctx.bindBuffer(ctx.ARRAY_BUFFER, h);

    const i = new Float32Array([-0.2, -0.9, 0, 0.4, -0.26, 0, 0, 0.7321, 0]);

    ctx.bufferData(ctx.ARRAY_BUFFER, i, ctx.STATIC_DRAW), (h.itemSize = 3), (h.numItems = 3);

    const j = ctx.createProgram();
    const k = ctx.createShader(ctx.VERTEX_SHADER);

    ctx.shaderSource(k, f);
    ctx.compileShader(k);

    const l = ctx.createShader(ctx.FRAGMENT_SHADER);

    ctx.shaderSource(l, g);
    ctx.compileShader(l);
    ctx.attachShader(j, k);
    ctx.attachShader(j, l);
    ctx.linkProgram(j);
    ctx.useProgram(j);

    j.vertexPosAttrib = ctx.getAttribLocation(j, 'attrVertex');
    j.offsetUniform = ctx.getUniformLocation(j, 'uniformOffset');

    ctx.enableVertexAttribArray(j.vertexPosArray);
    ctx.vertexAttribPointer(j.vertexPosAttrib, h.itemSize, ctx.FLOAT, !1, 0, 0);
    ctx.uniform2f(j.offsetUniform, 1, 1);
    ctx.drawArrays(ctx.TRIANGLE_STRIP, 0, h.numItems);

    const n = new Uint8Array(canvas.width * canvas.height * 4);
    ctx.readPixels(0, 0, canvas.width, canvas.height, ctx.RGBA, ctx.UNSIGNED_BYTE, n);

    const result = JSON.stringify(n).replace(/,?"[0-9]+":/g, '');

    if (debug) {
      document.body.appendChild(canvas);
    } else {
      ctx.clear(ctx.COLOR_BUFFER_BIT | ctx.DEPTH_BUFFER_BIT | ctx.STENCIL_BUFFER_BIT);
    }

    return murmurhash3_32_gc(result);
  } catch {
    return null;
  }
};

export const getWebglInfo = () => {
  try {
    const ctx = document.createElement('canvas').getContext('webgl');

    const result = {
      VERSION: String(ctx.getParameter(ctx.VERSION)),
      SHADING_LANGUAGE_VERSION: String(ctx.getParameter(ctx.SHADING_LANGUAGE_VERSION)),
      VENDOR: String(ctx.getParameter(ctx.VENDOR)),
      SUPORTED_EXTENSIONS: String(ctx.getSupportedExtensions()),
    };

    return result;
  } catch {
    return null;
  }
};

export const murmurhash3_32_gc = (key) => {
  const remainder = key.length & 3; // key.length % 4
  const bytes = key.length - remainder;
  const c1 = 0xcc9e2d51;
  const c2 = 0x1b873593;

  let h1, h1b, k1;

  for (let i = 0; i < bytes; i++) {
    k1 =
      (key.charCodeAt(i) & 0xff) |
      ((key.charCodeAt(++i) & 0xff) << 8) |
      ((key.charCodeAt(++i) & 0xff) << 16) |
      ((key.charCodeAt(++i) & 0xff) << 24);
    ++i;

    k1 = ((k1 & 0xffff) * c1 + ((((k1 >>> 16) * c1) & 0xffff) << 16)) & 0xffffffff;
    k1 = (k1 << 15) | (k1 >>> 17);
    k1 = ((k1 & 0xffff) * c2 + ((((k1 >>> 16) * c2) & 0xffff) << 16)) & 0xffffffff;

    h1 ^= k1;
    h1 = (h1 << 13) | (h1 >>> 19);
    h1b = ((h1 & 0xffff) * 5 + ((((h1 >>> 16) * 5) & 0xffff) << 16)) & 0xffffffff;
    h1 = (h1b & 0xffff) + 0x6b64 + ((((h1b >>> 16) + 0xe654) & 0xffff) << 16);
  }

  const i = bytes - 1;

  k1 = 0;

  switch (remainder) {
    case 3: {
      k1 ^= (key.charCodeAt(i + 2) & 0xff) << 16;
      break;
    }
    case 2: {
      k1 ^= (key.charCodeAt(i + 1) & 0xff) << 8;
      break;
    }
    case 1: {
      k1 ^= key.charCodeAt(i) & 0xff;
      break;
    }
  }

  k1 = ((k1 & 0xffff) * c1 + ((((k1 >>> 16) * c1) & 0xffff) << 16)) & 0xffffffff;
  k1 = (k1 << 15) | (k1 >>> 17);
  k1 = ((k1 & 0xffff) * c2 + ((((k1 >>> 16) * c2) & 0xffff) << 16)) & 0xffffffff;
  h1 ^= k1;

  h1 ^= key.length;

  h1 ^= h1 >>> 16;
  h1 = ((h1 & 0xffff) * 0x85ebca6b + ((((h1 >>> 16) * 0x85ebca6b) & 0xffff) << 16)) & 0xffffffff;
  h1 ^= h1 >>> 13;
  h1 = ((h1 & 0xffff) * 0xc2b2ae35 + ((((h1 >>> 16) * 0xc2b2ae35) & 0xffff) << 16)) & 0xffffffff;
  h1 ^= h1 >>> 16;

  return h1 >>> 0;
};

export const sha256 = (key) => {
  function rightRotate(value, amount) {
    return (value >>> amount) | (value << (32 - amount));
  };

  var mathPow = Math.pow;
  var maxWord = mathPow(2, 32);
  var lengthProperty = 'length'
  var i, j; // Used as a counter across the whole file
  var result = ''

  var words = [];
  var asciiBitLength = key[lengthProperty] * 8;

  //* caching results is optional - remove/add slash from front of this line to toggle
  // Initial hash value: first 32 bits of the fractional parts of the square roots of the first 8 primes
  // (we actually calculate the first 64, but extra values are just ignored)
  var hash = sha256.h = sha256.h || [];
  // Round constants: first 32 bits of the fractional parts of the cube roots of the first 64 primes
  var k = sha256.k = sha256.k || [];
  var primeCounter = k[lengthProperty];
  /*/
  var hash = [], k = [];
  var primeCounter = 0;
  //*/

  var isComposite = {};
  for (var candidate = 2; primeCounter < 64; candidate++) {
    if (!isComposite[candidate]) {
      for (i = 0; i < 313; i += candidate) {
        isComposite[i] = candidate;
      }
      hash[primeCounter] = (mathPow(candidate, .5) * maxWord) | 0;
      k[primeCounter++] = (mathPow(candidate, 1 / 3) * maxWord) | 0;
    }
  }

  key += '\x80' // Append Ƈ' bit (plus zero padding)
  while (key[lengthProperty] % 64 - 56) key += '\x00' // More zero padding
  for (i = 0; i < key[lengthProperty]; i++) {
    j = key.charCodeAt(i);
    if (j >> 8) return; // ASCII check: only accept characters in range 0-255
    words[i >> 2] |= j << ((3 - i) % 4) * 8;
  }
  words[words[lengthProperty]] = ((asciiBitLength / maxWord) | 0);
  words[words[lengthProperty]] = (asciiBitLength)

  // process each chunk
  for (j = 0; j < words[lengthProperty];) {
    var w = words.slice(j, j += 16); // The message is expanded into 64 words as part of the iteration
    var oldHash = hash;
    // This is now the undefinedworking hash", often labelled as variables a...g
    // (we have to truncate as well, otherwise extra entries at the end accumulate
    hash = hash.slice(0, 8);

    for (i = 0; i < 64; i++) {
      var i2 = i + j;
      // Expand the message into 64 words
      // Used below if 
      var w15 = w[i - 15], w2 = w[i - 2];

      // Iterate
      var a = hash[0], e = hash[4];
      var temp1 = hash[7]
        + (rightRotate(e, 6) ^ rightRotate(e, 11) ^ rightRotate(e, 25)) // S1
        + ((e & hash[5]) ^ ((~e) & hash[6])) // ch
        + k[i]
        // Expand the message schedule if needed
        + (w[i] = (i < 16) ? w[i] : (
          w[i - 16]
          + (rightRotate(w15, 7) ^ rightRotate(w15, 18) ^ (w15 >>> 3)) // s0
          + w[i - 7]
          + (rightRotate(w2, 17) ^ rightRotate(w2, 19) ^ (w2 >>> 10)) // s1
        ) | 0
        );
      // This is only used once, so *could* be moved below, but it only saves 4 bytes and makes things unreadble
      var temp2 = (rightRotate(a, 2) ^ rightRotate(a, 13) ^ rightRotate(a, 22)) // S0
        + ((a & hash[1]) ^ (a & hash[2]) ^ (hash[1] & hash[2])); // maj

      hash = [(temp1 + temp2) | 0].concat(hash); // We don't bother trimming off the extra ones, they're harmless as long as we're truncating when we do the slice()
      hash[4] = (hash[4] + temp1) | 0;
    }

    for (i = 0; i < 8; i++) {
      hash[i] = (hash[i] + oldHash[i]) | 0;
    }
  }

  for (i = 0; i < 8; i++) {
    for (j = 3; j + 1; j--) {
      var b = (hash[i] >> (j * 8)) & 255;
      result += ((b < 16) ? 0 : '') + b.toString(16);
    }
  }
  return result;
};

if (typeof window !== 'undefined') {
  window.getBrowserFingerprint = getBrowserFingerprint;
}

export default getBrowserFingerprint;
