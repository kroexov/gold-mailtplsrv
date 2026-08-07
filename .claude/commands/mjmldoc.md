# MJML Documentation Reference

Полный справочник по MJML — фреймворку для создания responsive HTML email.

## Структура документа

```mjml
<mjml>
  <mj-head>
    <!-- Метаданные, стили, атрибуты -->
  </mj-head>
  <mj-body>
    <!-- Контент email -->
  </mj-body>
</mjml>
```

### `<mjml>`
Корневой элемент.

| Атрибут | Значения | По умолчанию | Описание |
|---------|----------|-------------|----------|
| `owa` | `desktop` | — | Outlook Web App таргетинг |
| `lang` | string | — | Язык документа (`ru`, `en`) |
| `dir` | `ltr`, `rtl` | `ltr` | Направление текста |

---

## HEAD компоненты

### `<mj-attributes>`
Определяет дефолтные стили и переиспользуемые классы.

```mjml
<mj-attributes>
  <mj-all font-family="Arial, sans-serif" />        <!-- Глобально для всех -->
  <mj-text font-size="16px" line-height="24px" />    <!-- Дефолт для mj-text -->
  <mj-class name="blue" color="#0000ff" />            <!-- Именованный класс -->
</mj-attributes>
```

**Порядок разрешения атрибутов:**
1. Inline атрибуты на элементе
2. `mj-class` определения
3. Компонент-специфичные дефолты (`<mj-text>` в `<mj-attributes>`)
4. `<mj-all>` глобальные дефолты
5. MJML встроенные дефолты

### `<mj-class>`
Именованная группа атрибутов. Определяется внутри `<mj-attributes>`.

```mjml
<mj-attributes>
  <mj-class name="title" font-size="24px" font-weight="700" color="#000" />
  <mj-class name="btn-primary" background-color="#007bff" color="#fff" border-radius="4px" />
</mj-attributes>

<!-- Использование -->
<mj-text mj-class="title" padding="0">Заголовок</mj-text>
<mj-button mj-class="btn-primary" href="#">Кнопка</mj-button>
```

### `<mj-breakpoint>`
Кастомный responsive breakpoint (точка перехода desktop → mobile).

| Атрибут | Значения | По умолчанию |
|---------|----------|-------------|
| `width` | px | `480px` |

```mjml
<mj-breakpoint width="600px" />
```

### `<mj-font>`
Подключение веб-шрифтов.

| Атрибут | Описание |
|---------|----------|
| `name` | Имя шрифта |
| `href` | URL CSS (Google Fonts и т.д.) |

```mjml
<mj-font name="Roboto" href="https://fonts.googleapis.com/css?family=Roboto:300,400,500" />
```

### `<mj-preview>`
Текст превью в inbox (preheader). Только контент, без атрибутов.

```mjml
<mj-preview>Этот текст виден в списке писем до открытия</mj-preview>
```

### `<mj-title>`
Заголовок документа (`<title>` в HTML).

```mjml
<mj-title>Добро пожаловать!</mj-title>
```

### `<mj-style>`
CSS стили для email.

| Атрибут | Значения | Описание |
|---------|----------|----------|
| `inline` | `inline` | Инлайнить CSS в style="" атрибуты элементов |

**Два режима:**

```mjml
<!-- 1. Inline: CSS → style="" атрибуты на элементах -->
<!-- ✅ Надёжно в email-клиентах -->
<!-- ❌ Media queries НЕ работают (бейкаются в inline) -->
<mj-style inline="inline">
  .my-class { color: red; font-size: 16px; }
</mj-style>

<!-- 2. Non-inline: CSS → <style> в <head> -->
<!-- ✅ Media queries работают -->
<!-- ❌ Не все email-клиенты поддерживают <style> в head -->
<mj-style>
  .my-class { padding: 20px; }
  @media only screen and (max-width: 479px) {
    .my-class { padding: 10px !important; }
  }
</mj-style>
```

**Важно:** `css-class` на MJML компоненте применяется к **самому внешнему** HTML тегу. Для таргетинга внутренних элементов нужно смотреть сгенерированный HTML и использовать CSS-селекторы потомков.

### `<mj-include>`
Подключение внешних файлов.

| Атрибут | Значения | Описание |
|---------|----------|----------|
| `path` | string | Путь к файлу |
| `type` | `css`, `html` | Тип файла |
| `css-inline` | `inline` | Инлайнить CSS (для .css файлов) |

```mjml
<mj-include path="./header.mjml" />
<mj-include path="./styles.css" type="css" css-inline="inline" />
<mj-include path="./tracking.html" type="html" />
```

### `<mj-html-attributes>`
Inject кастомных HTML-атрибутов через CSS-селекторы.

```mjml
<mj-html-attributes>
  <mj-selector path=".custom-link">
    <mj-html-attribute name="data-tracking">true</mj-html-attribute>
  </mj-selector>
</mj-html-attributes>
```

---

## BODY — Layout компоненты

### `<mj-body>`
Корневой контейнер контента.

| Атрибут | Значения | По умолчанию | Описание |
|---------|----------|-------------|----------|
| `width` | px | `600px` | Ширина email |
| `background-color` | CSS color | — | Фон body |
| `css-class` | string | — | CSS класс |

```mjml
<mj-body width="648px" background-color="#f3f3f3">
  ...
</mj-body>
```

### `<mj-section>`
Строка layout — контейнер для колонок. На mobile колонки **стекаются вертикально**.

| Атрибут | Значения | По умолчанию | Описание |
|---------|----------|-------------|----------|
| `background-color` | CSS color | — | Фон секции |
| `background-url` | URL | — | Фоновое изображение |
| `background-size` | auto, cover, contain, px, % | `auto` | Размер фона |
| `background-position` | CSS position | `top center` | Позиция фона |
| `background-repeat` | repeat, no-repeat | `repeat` | Повтор фона |
| `border` | CSS border | — | Бордер |
| `border-radius` | px, % | — | Скругление |
| `direction` | ltr, rtl | `ltr` | Порядок колонок |
| `full-width` | full-width | — | 100% ширина (за пределы body width) |
| `padding` | px, % (1-4 значения) | `20px 0` | Внутренний отступ |
| `text-align` | left, center, right | `center` | Выравнивание |
| `css-class` | string | — | CSS класс |

```mjml
<mj-section padding="20px" background-color="#ffffff">
  <mj-column>...</mj-column>
  <mj-column>...</mj-column>
</mj-section>
```

### `<mj-column>`
Колонка внутри section. По умолчанию стекается на mobile (< breakpoint).

| Атрибут | Значения | По умолчанию | Описание |
|---------|----------|-------------|----------|
| `width` | px, % | авто (равные доли) | Ширина колонки |
| `background-color` | CSS color | — | Фон колонки |
| `inner-background-color` | CSS color | — | Фон контента (требует padding) |
| `padding` | px, % | — | Внешний отступ |
| `border` | CSS border | — | Бордер |
| `border-radius` | px, % | — | Скругление |
| `inner-border` | CSS border | — | Внутренний бордер |
| `inner-border-radius` | px, % | — | Внутреннее скругление |
| `direction` | ltr, rtl | `ltr` | Порядок дочерних элементов |
| `vertical-align` | top, middle, bottom | `top` | Вертикальное выравнивание |
| `css-class` | string | — | CSS класс |

**Ограничение:** сумма ширин колонок не может превышать ширину parent / 100%.

### `<mj-group>`
**Предотвращает стекание колонок на mobile.** Колонки внутри group сохраняют пропорции на всех экранах.

| Атрибут | Значения | По умолчанию | Описание |
|---------|----------|-------------|----------|
| `width` | px, % | авто | Ширина группы |
| `background-color` | CSS color | — | Фон |
| `direction` | ltr, rtl | `ltr` | Направление |
| `vertical-align` | top, middle, bottom | `top` | Выравнивание |
| `css-class` | string | — | CSS класс |

**Важно:** колонки в group — только `%` ширина, не `px`.

```mjml
<mj-section>
  <mj-group>
    <mj-column width="50%">Всегда слева</mj-column>
    <mj-column width="50%">Всегда справа</mj-column>
  </mj-group>
</mj-section>
```

### `<mj-wrapper>`
Группирует несколько `mj-section` с общим фоном, бордерами или padding.

| Атрибут | Значения | По умолчанию | Описание |
|---------|----------|-------------|----------|
| `background-color` | CSS color | — | Общий фон для секций |
| `background-url` | URL | — | Фоновое изображение |
| `background-size` | auto, cover, contain, px, % | `auto` | Размер фона |
| `background-position` | CSS position | `top center` | Позиция фона |
| `background-repeat` | repeat, no-repeat | `repeat` | Повтор фона |
| `border` | CSS border | — | Бордер |
| `border-radius` | px, % | — | Скругление |
| `full-width` | full-width | — | 100% ширина |
| `padding` | px, % (1-4 значения) | `20px 0` | Padding wrapper'а |
| `gap` | px | — | Вертикальный spacing между дочерними секциями |
| `direction` | ltr, rtl | `ltr` | Направление |
| `text-align` | left, center, right | `center` | Выравнивание |
| `css-class` | string | — | CSS класс |

**Отличие от mj-section:** wrapper содержит секции; секция содержит колонки. Wrapper позволяет задать единый фон/бордер для нескольких строк.

**Caveat:** При `full-width` на wrapper, дочерние секции с `full-width` возвращаются к стандартной ширине.

```mjml
<mj-wrapper padding="20px" background-color="#eee" gap="8px">
  <mj-section background-color="#fff">...</mj-section>
  <mj-section background-color="#fff">...</mj-section>
</mj-wrapper>
```

---

## BODY — Content компоненты

### `<mj-text>`
Текст и HTML контент. **Закрывающий тег обязателен.** Внутри — HTML, не MJML.

| Атрибут | Значения | По умолчанию | Описание |
|---------|----------|-------------|----------|
| `color` | CSS color | `#000000` | Цвет текста |
| `font-family` | string | `Ubuntu, Helvetica, Arial, sans-serif` | Шрифт |
| `font-size` | px | `13px` | Размер |
| `font-weight` | CSS weight | — | Жирность |
| `font-style` | CSS style | — | Курсив |
| `line-height` | px, %, unitless | `1` | Межстрочный |
| `letter-spacing` | px, em | — | Межбуквенный |
| `align` | left, center, right, justify | `left` | Выравнивание |
| `text-decoration` | underline, overline, none | — | Декорация |
| `text-transform` | capitalize, uppercase, lowercase, none | — | Трансформация |
| `padding` | px, % (1-4) | `10px 25px` | Отступы |
| `height` | px | — | Высота |
| `container-background-color` | CSS color | — | Фон контейнера |
| `css-class` | string | — | CSS класс |

### `<mj-image>`
Адаптивное изображение.

| Атрибут | Значения | По умолчанию | Описание |
|---------|----------|-------------|----------|
| `src` | URL | — | Источник |
| `width` | px | ширина колонки | Ширина |
| `height` | px, auto | `auto` | Высота |
| `max-height` | px, % | — | Макс. высота |
| `alt` | string | — | Alt текст |
| `align` | left, center, right | `center` | Выравнивание |
| `href` | URL | — | Ссылка по клику |
| `target` | string | `_blank` | Target ссылки |
| `border` | CSS border | — | Бордер |
| `border-radius` | px, % | — | Скругление |
| `padding` | px, % (1-4) | `10px 25px` | Отступы |
| `fluid-on-mobile` | boolean | — | 100% ширина на mobile |
| `srcset` | string | — | Responsive images |
| `sizes` | string | — | Sizes для srcset |
| `title` | string | — | Tooltip |
| `container-background-color` | CSS color | — | Фон контейнера |
| `css-class` | string | — | CSS класс |

### `<mj-button>`
Кнопка-ссылка. **Закрывающий тег обязателен.**

| Атрибут | Значения | По умолчанию | Описание |
|---------|----------|-------------|----------|
| `href` | URL | — | Ссылка |
| `background-color` | CSS color | `#414141` | Фон кнопки |
| `color` | CSS color | `#ffffff` | Цвет текста |
| `font-size` | px | `13px` | Размер шрифта |
| `font-family` | string | Ubuntu, Helvetica... | Шрифт |
| `font-weight` | CSS weight | — | Жирность |
| `font-style` | CSS style | — | Курсив |
| `border` | CSS border | — | Бордер |
| `border-radius` | px | `3px` | Скругление |
| `align` | left, center, right | `center` | Выравнивание кнопки |
| `vertical-align` | top, middle, bottom | `middle` | Верт. выравнивание |
| `width` | px, % | — | Ширина |
| `height` | px | — | Высота |
| `padding` | px, % (1-4) | `10px 25px` | Внешний отступ |
| `inner-padding` | px, % (1-4) | `10px 25px` | Внутренний отступ кнопки |
| `line-height` | px, % | `120%` | Межстрочный |
| `text-align` | left, center, right | — | Выравнивание текста |
| `text-decoration` | underline, overline, none | — | Декорация |
| `text-transform` | capitalize, uppercase, lowercase, none | — | Трансформация |
| `letter-spacing` | px, em | — | Межбуквенный |
| `target` | string | `_blank` | Target ссылки |
| `rel` | string | — | Rel атрибут |
| `title` | string | — | Tooltip |
| `container-background-color` | CSS color | — | Фон контейнера |
| `css-class` | string | — | CSS класс |

### `<mj-divider>`
Горизонтальная линия-разделитель.

| Атрибут | Значения | По умолчанию | Описание |
|---------|----------|-------------|----------|
| `border-width` | px | `4px` | Толщина |
| `border-style` | solid, dashed, dotted | `solid` | Стиль |
| `border-color` | CSS color | `#000000` | Цвет |
| `width` | px, % | `100%` | Ширина |
| `align` | left, center, right | `center` | Выравнивание |
| `padding` | px, % (1-4) | `10px 25px` | Отступы |
| `container-background-color` | CSS color | — | Фон контейнера |
| `css-class` | string | — | CSS класс |

### `<mj-spacer>`
Пустое пространство.

| Атрибут | Значения | По умолчанию | Описание |
|---------|----------|-------------|----------|
| `height` | px, % | `0px` | Высота |
| `padding` | px, % (1-4) | — | Отступы |
| `container-background-color` | CSS color | — | Фон |
| `css-class` | string | — | CSS класс |

### `<mj-table>`
HTML таблица. **Закрывающий тег.** Внутри — стандартный HTML `<tr>/<td>`.

| Атрибут | Значения | По умолчанию | Описание |
|---------|----------|-------------|----------|
| `width` | px, %, auto | `100%` | Ширина |
| `align` | left, center, right | `left` | Выравнивание |
| `border` | CSS border | `none` | Бордер |
| `cellpadding` | integer | `0` | Cellpadding |
| `cellspacing` | integer | `0` | Cellspacing |
| `color` | CSS color | `#000000` | Цвет текста |
| `font-family` | string | — | Шрифт |
| `font-size` | px | `13px` | Размер |
| `line-height` | px, % | `22px` | Межстрочный |
| `padding` | px, % (1-4) | `10px 25px` | Отступы |
| `table-layout` | auto, fixed | `auto` | Layout алгоритм |
| `role` | none, presentation | — | ARIA role |
| `container-background-color` | CSS color | — | Фон контейнера |
| `css-class` | string | — | CSS класс |

### `<mj-raw>`
Сырой HTML без обработки MJML. **Закрывающий тег.**

| Атрибут | Значения | Описание |
|---------|----------|----------|
| `position` | `file-start` | Вставить перед `<!DOCTYPE>` |

```mjml
<mj-raw>
  <!-- htmlmin:ignore -->
  {{#if condition}}
    <div>Template content</div>
  {{/if}}
  <!-- htmlmin:ignore -->
</mj-raw>
```

`<!-- htmlmin:ignore -->` защищает template syntax от минификации.

---

## Интерактивные компоненты

### `<mj-accordion>`
Раскрывающиеся секции (поддержка ограничена email-клиентами).

| Атрибут | Значения | По умолчанию |
|---------|----------|-------------|
| `border` | CSS border | `2px solid black` |
| `font-family` | string | Ubuntu, Helvetica... |
| `icon-position` | left, right | `right` |
| `icon-width` | px, % | `32px` |
| `icon-height` | px, % | `32px` |
| `icon-align` | top, middle, bottom | `middle` |
| `icon-wrapped-url` | URL | стрелка вниз |
| `icon-unwrapped-url` | URL | стрелка вверх |
| `padding` | px, % | `10px 25px` |
| `css-class` | string | — |

Дочерние: `<mj-accordion-element>` → `<mj-accordion-title>` + `<mj-accordion-text>`

### `<mj-carousel>`
Карусель изображений.

| Атрибут | Значения | По умолчанию |
|---------|----------|-------------|
| `align` | left, center, right | `center` |
| `border-radius` | px, % | `6px` |
| `icon-width` | px, % | `44px` |
| `left-icon` | URL | стрелка влево |
| `right-icon` | URL | стрелка вправо |
| `thumbnails` | visible, hidden | `hidden` |
| `tb-border` | CSS border | `2px solid transparent` |
| `tb-selected-border-color` | CSS color | `#ccc` |
| `css-class` | string | — |

Дочерние: `<mj-carousel-image src="" alt="" href="" />`

### `<mj-navbar>`
Навигационное меню с hamburger на mobile.

| Атрибут | Значения | По умолчанию |
|---------|----------|-------------|
| `align` | left, center, right | `center` |
| `hamburger` | `hamburger` | — |
| `ico-color` | CSS color | `#000000` |
| `ico-font-size` | px, % | `30px` |
| `ico-open` | ASCII | `&#9776;` |
| `ico-close` | ASCII | `&#8855;` |
| `css-class` | string | — |

Дочерние: `<mj-navbar-link href="" color="" font-size="">Text</mj-navbar-link>`

**Hamburger работает только в Apple Mail.**

### `<mj-social>`
Блок соцсетей с иконками.

| Атрибут | Значения | По умолчанию |
|---------|----------|-------------|
| `mode` | horizontal, vertical | `horizontal` |
| `align` | left, center, right | `center` |
| `icon-size` | px, % | `20px` |
| `border-radius` | px, % | `3px` |
| `padding` | px, % | `10px 25px` |
| `css-class` | string | — |

Дочерние: `<mj-social-element name="facebook" href="...">Facebook</mj-social-element>`

Встроенные пресеты `name`: facebook, twitter, x, google, pinterest, linkedin, tumblr, xing, github, instagram, web, snapchat, youtube, vimeo, medium, soundcloud, dribbble.

### `<mj-hero>`
Полноширинный баннер с фоновым изображением.

| Атрибут | Значения | По умолчанию |
|---------|----------|-------------|
| `mode` | fixed-height, fluid-height | `fluid-height` |
| `height` | px, % | `0px` |
| `background-url` | URL | — |
| `background-width` | px, % | — |
| `background-height` | px, % | — |
| `background-color` | CSS color | `#ffffff` |
| `background-position` | CSS position | `center center` |
| `vertical-align` | top, middle, bottom | `top` |
| `padding` | px, % | `0px` |
| `css-class` | string | — |

---

## Responsive design

### Автоматическое поведение
- **Колонки стекаются** на `< breakpoint` (480px по умолчанию)
- **Ширина** колонок становится 100%
- **Изображения** масштабируются до ширины колонки

### Предотвращение стекания
```mjml
<mj-group>
  <mj-column width="50%">Не стекается</mj-column>
  <mj-column width="50%">Не стекается</mj-column>
</mj-group>
```

### Кастомный breakpoint
```mjml
<mj-breakpoint width="600px" />
```

### Media queries
Только в non-inline `<mj-style>`:
```mjml
<mj-style>
  @media only screen and (max-width: 479px) {
    .mobile-hidden { display: none !important; }
    .mobile-full { width: 100% !important; }
  }
</mj-style>
```

### fluid-on-mobile
```mjml
<mj-image src="..." fluid-on-mobile="true" />
```

### full-width секции
```mjml
<mj-section full-width="full-width" background-color="#eee">
  <!-- Фон на 100% ширины viewport, контент ограничен body width -->
</mj-section>
```

---

## Общие атрибуты

Большинство content-компонентов поддерживают:
- `padding` — внутренний отступ (px, %, 1-4 значения)
- `css-class` — CSS класс (применяется к внешнему HTML элементу)
- `mj-class` — ссылка на `<mj-class>` из `<mj-attributes>`
- `container-background-color` — фон контейнера элемента

## Ссылки

- Документация: https://documentation.mjml.io/
- GitHub: https://github.com/mjmlio/mjml
- Playground: https://mjml.io/try-it-live
