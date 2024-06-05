/**
 * @param {string} a 
 * @param {string} b
 * @returns {string}
 */
function sum(a, b) {
  if (a.length > b.length) {
    let temp = a;
    a = b
    b = temp
  }
  let result = new Array(b.length + 3)
  a = a.split('').reverse().join('')
  b = b.split('').reverse().join('')

  let add = 0;
  for (let i = 0; i < a.length; i++) {
    let sum = add + +a[i] + +b[i]
    if (sum > 9) {
      add = 1
      sum = sum % 10
    } else {
      add = 0
    }
    result.unshift(sum)
  }

  for (let i = a.length; i < b.length; i++) {
    let sum = +b[i] + add
    if (sum > 9) {
      add = 1
      sum = sum % 10
    } else {
      add = 0
    }
    result.unshift(sum)
  }

  if (add) {
    result.unshift(add)
  }

  return result
}
