//------------------------------------------------------------------------------------------------//
//--  Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)   --//
//------------------------------------------------------------------------------------------------//
import { LitElement, css, html } from "lit";
import { classMap } from "lit/directives/class-map.js";

export class KolobConversationList extends LitElement {
  static properties = {
    _conversations: { state: true },
    _selected: { state: true },
  };

  constructor() {
    super();
    this._conversations = [
      {
        id: "first",
        name: "First Conversation",
      },
      {
        id: "second",
        name: "Second Conversation",
      },
      {
        id: "third",
        name: "Third Conversation",
      },
    ];
    this._selected = "";
  }

  render() {
    return html`
      <h3>Conversations</h3>
      ${this._conversations.map(
        (c) => html`
          <div class=${classMap({ selected: c.id === this._selected })}>
            <span data-id=${c.id} data-name=${c.name} @mouseup=${this._handleConversationMouseUp}>
              ${c.name}
            </span>
          </div>
        `,
      )}
    `;
  }

  _handleConversationMouseUp(e) {
    if (this._selected == e.target.dataset.id) {
      return;
    }

    this._selected = e.target.dataset.id;
    const detail = { id: this._selected, name: e.target.dataset.name };
    const event = new CustomEvent("kolob-conversation-selected", { bubbles: true, detail });
    this.dispatchEvent(event);
  }
}

LitElement.styles = css`
  :host {
    display: flex;
    flex-direction: column;
  }

  h3 {
    margin: 0px 4px 4px 4px;
  }

  div:hover {
    background-color: #cccccc;
    cursor: pointer;
  }

  div.selected {
    background-color: #cccccc;
  }

  span {
    display: block;
    padding: 4px 12px;
  }
`;
