//------------------------------------------------------------------------------------------------//
//--  Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)   --//
//------------------------------------------------------------------------------------------------//
import { LitElement, css, html } from "lit";
import { createRef, ref } from "lit/directives/ref.js";
import { classMap } from "lit/directives/class-map.js";

import { define } from "../define";

import { KolobUserContextConsumer } from "./context/KolobUserContext";

import { KolobMessage } from "./KolobMessage";
import { KolobMessageWriter } from "./KolobMessageWriter";

export class KolobMessageArea extends LitElement {
  static properties = {
    conversation: { attribute: false },
    _messages: { state: true },
  };

  constructor() {
    super();
    this.conversation = null;
    this._userContext = new KolobUserContextConsumer(this);
    this._messages = {
      first: [
        {
          author: {
            username: "jdoe",
            name: "John Doe",
          },
          created: new Date(),
          updated: null,
          content: "Hello there!",
        },
      ],
      second: [],
      third: [],
    };
    this._messageContainerRef = createRef();
  }

  updated() {
    const ms = this._messageContainerRef.value;
    if (!ms) {
      return;
    }
    const shouldScroll = ms.scrollTop + ms.clientHeight !== ms.scrollHeight;
    if (shouldScroll) {
      ms.scrollTop = ms.scrollHeight;
    }
  }

  render() {
    if (this.conversation === null) {
      return html`
        <div class="none-selected">
          <div>Please select a conversation to view messages</div>
        </div>
      `;
    }

    return html`
      <h3>${this.conversation.name}</h3>
      <div class="scroll-box-outer" ${ref(this._messageContainerRef)}>
        <div class="scroll-box-inner">${this._renderMessages()}</div>
      </div>
      <div class="writer">
        <kolob-message-writer
          @kolob-message-ready=${this._handleMessageReady}
        ></kolob-message-writer>
      </div>
    `;
  }

  _renderMessages() {
    const messages = this._messages[this.conversation.id];
    return messages.map((m, i) => {
      const right = m.author.username === this._userContext.value.username;
      const left = !right;
      const classes = { message: true, left, right };
      const isSameUserAsLast = messages[i - 1]?.author.username != m.author.username;
      const author = isSameUserAsLast ? m.author.name : "";
      return html`<kolob-message
        class=${classMap(classes)}
        author=${author}
        content=${m.content}
      ></kolob-message>`;
    });
  }

  _handleMessageReady(e) {
    const d = new Date();
    this._messages[this.conversation.id] = [
      ...this._messages[this.conversation.id],
      {
        author: this._userContext.value,
        created: d,
        updated: d,
        content: e.detail,
      },
    ];
    this.requestUpdate();
  }
}

KolobMessageArea.styles = css`
  :host {
    flex: 1;
    align-self: stretch;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  div.none-selected {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  h3 {
    margin: 0px 4px 8px 4px;
  }

  div.scroll-box-outer {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: auto;
  }

  div.scroll-box-inner {
    display: flex;
    flex-direction: column;
    min-height: min-content;
  }

  div.writer {
    display: flex;
  }
`;

define("kolob-message", KolobMessage);
define("kolob-message-writer", KolobMessageWriter);
