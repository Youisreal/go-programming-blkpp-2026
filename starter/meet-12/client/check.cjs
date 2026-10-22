// Jalankan: node client/check.cjs
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const nodes = new Map();
const node = (id) => {
  if (!nodes.has(id)) nodes.set(id, {});
  return nodes.get(id);
};
node('item-template').content = { cloneNode: () => ({ querySelector: () => ({}) }) };
node('items').append = () => {};
const context = vm.createContext({
  document: { getElementById: node }, Intl, AbortSignal,
  fetch: async (path, options) => {
    assert.equal(path, 'http://localhost:8080/api/private/orders');
    assert.equal(options.method, 'POST');
    assert.equal(options.headers.Authorization, 'Bearer test-token');
    assert.deepEqual(JSON.parse(options.body), { items: [{ product_id: 7, quantity: 2 }] });
    return { ok: true, json: async () => ({ data: { id: 42 } }) };
  }
});
const html = fs.readFileSync(`${__dirname}/index.html`, 'utf8');
vm.runInContext(html.match(/<script>([\s\S]*?)<\/script>/)[1], context);
const item = (id, quantity) => ({ querySelector: (selector) => ({ value: selector.includes('product_id') ? id : quantity }) });
node('items').children = [item('7', '2')];
assert.equal(vm.runInContext('JSON.stringify(readItems())', context), '[{"product_id":7,"quantity":2}]');
node('items').children.push(item('7', '1'));
assert.throws(() => vm.runInContext('readItems()', context), /duplikat/);
node('items').children = [item('0', '2')];
assert.throws(() => vm.runInContext('readItems()', context), /positif/);
vm.runInContext("token = 'test-token'; request('/api/private/orders', { items: [{ product_id: 7, quantity: 2 }] })", context)
  .then((order) => { assert.equal(order.id, 42); console.log('Validasi item dan request order: OK'); })
  .catch((error) => { console.error(error); process.exitCode = 1; });
