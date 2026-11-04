import {createTippy} from '../../modules/tippy.ts';
import {GET} from '../../modules/fetch.ts';
import {showErrorToast} from '../../modules/toast.ts';
import {replaceTextareaSelection} from './EditorMarkdown.ts';

export type SavedReply = {
  id: number,
  title: string,
  content: string,
};

export function filterSavedReplies(replies: Array<SavedReply>, query: string): Array<SavedReply> {
  const needle = query.trim().toLowerCase();
  if (!needle) return replies;
  return replies.filter((reply) => reply.title.toLowerCase().includes(needle));
}

export function renderSavedReplies(list: HTMLElement, replies: Array<SavedReply>, emptyText: string) {
  list.replaceChildren(...replies.map((reply) => {
    const item = document.createElement('button');
    item.type = 'button';
    item.className = 'item';
    item.textContent = reply.title;
    item.setAttribute('data-saved-reply-id', String(reply.id));
    return item;
  }));
  if (!replies.length) {
    const empty = document.createElement('div');
    empty.className = 'tw-p-2';
    empty.textContent = emptyText;
    list.append(empty);
  }
}

let cachedReplies: Array<SavedReply> | null = null; // shared by all editors of the page, the list can be several MB

export async function loadSavedReplies(url: string): Promise<Array<SavedReply>> {
  if (!cachedReplies) {
    const resp = await GET(url);
    if (!resp.ok) throw new Error(`unexpected status ${resp.status}`);
    cachedReplies = await resp.json();
  }
  return cachedReplies!;
}

export function initSavedRepliesButton(container: HTMLElement, textarea: HTMLTextAreaElement) {
  const button = container.querySelector('.markdown-button-saved-replies');
  if (!button) return;

  const panel = container.querySelector<HTMLElement>('.markdown-saved-replies-panel')!;
  const filter = panel.querySelector<HTMLInputElement>('.saved-replies-filter')!;
  const list = panel.querySelector<HTMLElement>('.saved-replies-list')!;
  const emptyText = panel.getAttribute('data-empty-text')!;
  let replies: Array<SavedReply> = [];

  // the button already owns a tippy for its tooltip, so the panel attaches to its group
  const panelTippy = createTippy(button.parentElement!, {
    content: panel,
    role: 'dialog', // the "menu" role's key handling would swallow space and enter in the filter input
    theme: 'menu',
    trigger: 'manual',
    placement: 'bottom-end',
    hideOnClick: true,
    interactive: true,
  });

  button.addEventListener('click', async () => {
    panelTippy.show();
    try {
      replies = await loadSavedReplies(panel.getAttribute('data-url')!);
    } catch {
      panelTippy.hide();
      showErrorToast('Failed to load saved replies');
      return;
    }
    renderSavedReplies(list, filterSavedReplies(replies, filter.value), emptyText);
    filter.focus();
  });

  filter.addEventListener('input', () => {
    renderSavedReplies(list, filterSavedReplies(replies, filter.value), emptyText);
  });

  list.addEventListener('click', (e) => {
    const item = (e.target as HTMLElement).closest('[data-saved-reply-id]');
    if (!item) return;
    const id = Number(item.getAttribute('data-saved-reply-id'));
    replaceTextareaSelection(textarea, replies.find((reply) => reply.id === id)!.content);
    panelTippy.hide();
  });
}
