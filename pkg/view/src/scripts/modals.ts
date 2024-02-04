export function openModel(selector: string) {
  document.querySelector<HTMLDialogElement>(selector)?.showModal()
}

export function closeModel(selector: string) {
  document.querySelector<HTMLDialogElement>(selector)?.close()
}