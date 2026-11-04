import {GET} from '../../modules/fetch.ts';
import {filterSavedReplies, loadSavedReplies, loadSavedReplyContent, renderSavedReplies} from './SavedReplies.ts';

vi.mock('../../modules/fetch.ts', () => ({GET: vi.fn()}));

const replies = [
  {id: 1, title: 'Duplicate'},
  {id: 2, title: 'Needs reproduction'},
];

test('filterSavedReplies matches titles case-insensitively', () => {
  expect(filterSavedReplies(replies, '  ')).toEqual(replies);
  expect(filterSavedReplies(replies, 'DUPL')).toEqual([replies[0]]);
  expect(filterSavedReplies(replies, 'closing')).toEqual([]);
});

test('renderSavedReplies escapes titles and shows the empty text', () => {
  const list = document.createElement('div');
  renderSavedReplies(list, [{id: 3, title: '<b>bold</b>'}], 'empty');
  const item = list.querySelector('button')!;
  expect(item.textContent).toBe('<b>bold</b>');
  expect(item.getAttribute('data-saved-reply-id')).toBe('3');

  renderSavedReplies(list, [], 'empty');
  expect(list.querySelector('button')).toBeNull();
  expect(list.textContent).toBe('empty');
});

test('loadSavedReplies caches a successful load only', async () => {
  vi.mocked(GET).mockResolvedValueOnce({ok: false, status: 500} as Response);
  await expect(loadSavedReplies('/replies')).rejects.toThrow();
  vi.mocked(GET).mockResolvedValueOnce({ok: true, json: () => Promise.resolve(replies)} as Response);
  expect(await loadSavedReplies('/replies')).toEqual(replies);
  expect(await loadSavedReplies('/replies')).toEqual(replies);
  expect(GET).toHaveBeenCalledTimes(2);
  expect(GET).toHaveBeenLastCalledWith('/replies/list');
});

test('loadSavedReplyContent fetches the picked reply', async () => {
  vi.mocked(GET).mockResolvedValueOnce({ok: true, json: () => Promise.resolve({id: 2, title: 'Needs reproduction', content: 'Steps?'})} as Response);
  expect(await loadSavedReplyContent('/replies', 2)).toBe('Steps?');
  expect(GET).toHaveBeenLastCalledWith('/replies/2');
  vi.mocked(GET).mockResolvedValueOnce({ok: false, status: 404} as Response);
  await expect(loadSavedReplyContent('/replies', 3)).rejects.toThrow();
});
