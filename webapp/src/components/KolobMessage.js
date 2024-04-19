//------------------------------------------------------------------------------------------------//
//--  Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)   --//
//------------------------------------------------------------------------------------------------//
import { LitElement, css, html } from "lit";

export class KolobMessage extends LitElement {
  static properties = {
    author: { type: String },
    content: { type: String },
  };

  constructor() {
    super();
    this.author = "";
    this.content = "";
  }

  render() {
    return html`
      <div class="message-container">
        <div class="author">${this.author}</div>
        <div class="message">${this.content}</div>
      </div>
    `;
  }
}

KolobMessage.styles = css`
  :host(.left) {
    align-self: flex-start;
  }

  :host(.right) {
    align-self: flex-end;
  }

  div.message-container {
    display: flex;
    flex-direction: column;
    margin: 2px 16px;
  }

  div.author {
    font-size: 0.8em;
    padding-left: 5px;
  }

  div.message {
    background-color: #cccccc;
    border-radius: 5px;
    padding: 4px 8px;
    position: relative;
    white-space: pre-wrap;
  }

  :host(.left) div.message:before {
    content: "";
    width: 0px;
    height: 0px;
    position: absolute;
    border-left: 4px solid transparent;
    border-right: 4px solid #cccccc;
    border-top: 4px solid transparent;
    border-bottom: 4px solid #cccccc;
    left: -5px;
    bottom: 0px;
  }

  :host(.right) div.message:before {
    content: "";
    width: 0px;
    height: 0px;
    position: absolute;
    border-left: 4px solid #cccccc;
    border-right: 4px solid transparent;
    border-top: 4px solid transparent;
    border-bottom: 4px solid #cccccc;
    right: -5px;
    bottom: 0px;
  }
`;
