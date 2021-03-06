import m  from 'mithril'
import {AddShare} from "./add_share";
import * as api from "./api";

export const Shares = (vnode) => {
    let { shares } = vnode.attrs;

    const reload = () => {
        api.getShares()
            .then((data) => shares = data)
            .then(m.redraw);
    };

    return {
        view: () => {
            shares = shares.sort((a, b) => a.id < b.id);

            return m('.shares-container',
                m('.shares-header', m('p.title', '↓ share your queue')),
                m('.shares', [
                    m(AddShare, {key: 'add', reload}),
                    ...shares.map((share) => m(Share, {key: share.id, share, reload}))
                ])
            );
        }
    };
};

export const Share = (vnode) => {
    const { share } = vnode.attrs;
    let enabled = share.enabled;
    let disabled = false;

    const setEnabled = () => {
        disabled = true;
        api.setEnabled(share.id, !enabled)
            .then(() => enabled = !enabled)
            .then(m.redraw);
    };

    return {
        view: () => m('.share.card',
            m('.info',
                m('img.image', {src: share.image_url}),
                m('.name', m('p', share.name)),
            ),
            m('.enabled',
                m('input', { id: `share-enabled-${share.id}`, type: 'checkbox', checked: enabled, onchange: setEnabled }),
                m('label', { for: `share-enabled-${share.id}` }, 'Enabled'),
            )
        )
    };
};
