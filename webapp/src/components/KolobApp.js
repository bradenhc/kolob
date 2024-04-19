//------------------------------------------------------------------------------------------------//
//--  Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)   --//
//------------------------------------------------------------------------------------------------//
import { LitElement, css, html } from "lit";

import { define } from "../define";

import { KolobUserContextProvider } from "./context/KolobUserContext";

import { KolobConversationList } from "./KolobConversationList";
import { KolobMessageArea } from "./KolobMessageArea";
import { KolobSignIn } from "./KolobSignIn";

export class KolobApp extends LitElement {
  static properties = {
    _selectedConversation: { state: true },
  };

  constructor() {
    super();
    this._userContext = new KolobUserContextProvider(this);
    this._selectedConversation = null;
  }

  render() {
    if (this._userContext.value === undefined) {
      return html`
        <kolob-sign-in @kolob-authenticated=${this._handleAuthenticated}></kolob-sign-in>
      `;
    }

    return html`
      <div class="group">
        <div class="conversations">
          <kolob-conversation-list
            @kolob-conversation-selected=${this._handleConversationSelected}
          ></kolob-conversation-list>
        </div>
        <div class="messages">
          <kolob-message-area .conversation=${this._selectedConversation}></kolob-messages-area>
        </div>
      </div>
    `;
  }

  _handleAuthenticated(e) {
    this._userContext.setValue(e.detail);
    this.requestUpdate();
  }

  _handleConversationSelected(e) {
    this._selectedConversation = e.detail;
  }
}

KolobApp.styles = css`
  :host {
    display: flex;
    flex-direction: column;
    width: 100%;
    height: 100%;
    overflow: hidden;
  }

  .group {
    flex: 1;
    align-self: stretch;
    display: flex;
    align-items: stretch;
    overflow: hidden;
  }

  .conversations {
    display: flex;
    flex-direction: column;
    align-items: stretch;
    border-right: 1px solid #cccccc;
  }

  .messages {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: stretch;
    overflow: hidden;
  }
`;

define("kolob-conversation-list", KolobConversationList);
define("kolob-message-area", KolobMessageArea);
define("kolob-sign-in", KolobSignIn);
