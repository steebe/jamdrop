import m from 'mithril';
import {Sharers} from "./sharers";
import {Shares} from "./shares";
import jam from '../img/jam.svg';
import * as api from "./api";
import toaster from "./toaster";

export const Main = (vnode) => {
    const { userData } = vnode.attrs;
    const user = userData.user;

    return {
        view: () => m('main',
            m('.welcome',
                m('.logo', m('img', {src: jam}), m('p', 'JamDrop')),
                toaster.message,
                m('.user',
                    m('p', `Hi, ${user.name} 👋`),
                    m(Settings, {user})
                )
            ),
            m(Sharers, {sharers: userData.sharers}),
            m(Shares, {shares: userData.shares}),
        )
    };
};

export const Settings = (vnode) => {
    const {user} = vnode.attrs;
    let stayActive = user.stay_active;
    let phoneNumber = user.phone_number;

    return {
        view: () => {
            let stayActiveDisabled = false;
            let phoneNumberDisabled = false;

            const setStayActive = () => {
                stayActiveDisabled = true;
                api.updateUser({ stay_active: !stayActive })
                    .then(() => stayActive = !stayActive)
                    .then(m.redraw);
            };

            const setPhoneNumber = (event) => {
                const newPhoneNumber = event.target.value;

                phoneNumberDisabled = true;
                api.updateUser({ phone_number: newPhoneNumber })
                    .then(() => phoneNumber = newPhoneNumber)
                    .then(m.redraw);
            };

            return m('.settings',
                // m('.setting',
                //     m('input', {
                //         type: 'text',
                //         placeholder: 'phone number',
                //         value: phoneNumber,
                //         onblur: setPhoneNumber,
                //     })
                // ),
                m('.setting',
                    m('input#stay-active', {
                        type: 'checkbox',
                        disabled: stayActiveDisabled,
                        checked: stayActive,
                        onchange: setStayActive
                    }),
                    m('label', {for: 'stay-active'}, 'Stay active')
                ),
            );
        }
    };
};
