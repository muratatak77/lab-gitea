import {GET} from '../../modules/fetch.ts';
import {filterSavedReplies, loadSavedReplies, renderSavedReplies} from './SavedReplies.ts';

vi.mock('../../modules/fetch.ts', () => ({GET: vi.fn()}));

const replies = [
  {id: 1, title: 'Duplicate', content: 'Closing as a duplicate.'},
  {id: 2, title: 'Needs reproduction', content: 'Steps?'},
];

test('filterSavedReplies matches titles case-insensitively', () => {
  expect(filterSavedReplies(replies, '  ')).toEqual(replies);
  expect(filterSavedReplies(replies, 'DUPL')).toEqual([replies[0]]);
  expect(filterSavedReplies(replies, 'closing')).toEqual([]);
});

test('renderSavedReplies escapes titles and shows the empty text', () => {
  const list = document.createElement('div');
  renderSavedReplies(list, [{id: 3, title: '<b>bold</b>', content: 'x'}], 'empty');
  const item = list.querySelector('button')!;
  expect(item.textContent).toBe('<b>bold</b>');
  expect(item.getAttribute('data-saved-reply-id')).toBe('3');

  renderSavedReplies(list, [], 'empty');
  expect(list.querySelector('button')).toBeNull();
  expect(list.textContent).toBe('empty');
});

test('loadSavedReplies caches a successful load only', async () => {
  vi.mocked(GET).mockResolvedValueOnce({ok: false, status: 500} as Response);
  await expect(loadSavedReplies('/list')).rejects.toThrow();
  vi.mocked(GET).mockResolvedValueOnce({ok: true, json: () => Promise.resolve(replies)} as Response);
  expect(await loadSavedReplies('/list')).toEqual(replies);
  expect(await loadSavedReplies('/list')).toEqual(replies);
  expect(GET).toHaveBeenCalledTimes(2);
});
