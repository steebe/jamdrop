import m from 'mithril';

const activityHeader = {
    name: 'name',
    dropped: 'dropped',
    when: 'when',
};

let activityTable = [{}];

export const Activity = () => {
    return {
        view: () => {
            return m('table',[
                m('thead', m('p.title','→ activity')),
                m('tr', [
                    m('td', m('p.message', activityHeader.name)),
                    m('td', m('p.message', activityHeader.dropped)),
                    m('td', m('p.message', activityHeader.when)),
                ]),
                m('tr', [
                    m('td', 'Hey!'),
                    m('td', 'You!'),
                ])
            ]);
        }
    }
}

export const log = (message) => {

}