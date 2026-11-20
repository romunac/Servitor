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

document.querySelector('#year').textContent = new Date().getFullYear();
