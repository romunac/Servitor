const menuToggle = document.querySelector('.menu-toggle');
const primaryNav = document.querySelector('#primary-nav');

menuToggle.addEventListener('click', () => {
  const expanded = menuToggle.getAttribute('aria-expanded') === 'true';
  menuToggle.setAttribute('aria-expanded', String(!expanded));
  menuToggle.setAttribute('aria-label', expanded ? 'Open navigation' : 'Close navigation');
  primaryNav.classList.toggle('is-open', !expanded);
});

primaryNav.addEventListener('click', (event) => {
  if (event.target.closest('a') && window.matchMedia('(max-width: 700px)').matches) {
    menuToggle.setAttribute('aria-expanded', 'false');
    menuToggle.setAttribute('aria-label', 'Open navigation');
    primaryNav.classList.remove('is-open');
  }
});

document.querySelectorAll('.copy-command').forEach((button) => {
  button.addEventListener('click', async () => {
    const code = button.closest('.code-window')?.querySelector('code')?.textContent;
    if (!code || !navigator.clipboard?.writeText) {
      button.textContent = 'Unavailable';
      return;
    }

    try {
      await navigator.clipboard.writeText(code);
      button.textContent = 'Copied';
      window.setTimeout(() => { button.textContent = 'Copy'; }, 1600);
    } catch {
      button.textContent = 'Copy failed';
    }
  });
});

document.querySelector('#year').textContent = new Date().getFullYear();
