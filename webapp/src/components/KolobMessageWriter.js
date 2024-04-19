//------------------------------------------------------------------------------------------------//
//--  Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)   --//
//------------------------------------------------------------------------------------------------//
import { LitElement, css, html } from "lit";

export class KolobMessageWriter extends LitElement {
  static properties = {
    _content: { state: true },
  };

  constructor() {
    super();
    this._content = "";
  }

  render() {
    return html`
      <div class="textarea-wrapper">
        <textarea
          rows="3"
          @change=${this._handleChange}
          @keydown=${this._handleKeyDown}
          @keyup=${this._handleKeyUp}
          .value=${this._content}
        ></textarea>
      </div>
      <div class="actions">
        <button @mouseup=${this._handleSendMouseUp}>Send</button>
      </div>
    `;
  }

  _handleChange(e) {
    this._content = e.target.value;
  }

  _handleKeyDown(e) {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      e.stopPropagation();
      this._send();
    }
  }

  _handleSendMouseUp(e) {
    this._send();
  }

  _handleKeyUp(e) {
    this._handleChange(e);
  }

  _send() {
    if (this._content === "") {
      return;
    }

    const event = new CustomEvent("kolob-message-ready", { bubbles: true, detail: this._content });
    this.dispatchEvent(event);
    this._content = "";
  }
}

KolobMessageWriter.styles = css`
  :host {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  div.textarea-wrapper {
    flex: 1;
    margin: 8px;
    padding: 4px;
    border: 1px solid #cccccc;
    border-radius: 5px;
  }

  textarea {
    width: 100%;
    border: none;
    resize: none;
  }

  textarea:focus-visible {
    outline: none;
  }

  div.actions {
    display: flex;
    flex-direction: row-reverse;
    margin: 0px 8px 8px 8px;
    padding: 0px 4px;
  }
`;
