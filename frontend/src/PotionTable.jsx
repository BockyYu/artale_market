import { HP_POTIONS, MP_POTIONS } from './potionData'

const hpRateColor = (rate, stat) => rate > 0.9 && stat >= 1000 ? '#c62828' : 'inherit'
const mpRateColor = (rate) => rate > 0.49 ? '#c62828' : 'inherit'

function PotionSection({ title, potions, statKey, statLabel, rateColor }) {
  return (
    <div className="potion-section">
      <h2 className="potion-title">{title}</h2>
      <p className="potion-subtitle">轉換率 = 補量 ÷ 價格，數值越高越划算</p>
      <div className="table-wrapper">
        <table className="potion-table">
          <thead>
            <tr>
              <th>藥水名</th>
              <th>補量 ({statLabel})</th>
              <th>價格</th>
              <th>轉換率</th>
              <th>購買地</th>
            </tr>
          </thead>
          <tbody>
            {potions.map((potion) => {
              const bestRate = Math.max(...potion.entries.map(e => e.rate))
              return potion.entries.map((entry, idx) => (
                <tr key={`${potion.name}-${idx}`} className={entry.rate === bestRate ? 'best-entry' : ''}>
                  {idx === 0 && (
                    <>
                      <td rowSpan={potion.entries.length} className="potion-name">
                        {potion.name}
                      </td>
                      <td rowSpan={potion.entries.length} className="potion-hp">
                        {potion[statKey].toLocaleString()}
                      </td>
                    </>
                  )}
                  <td className="potion-price">{entry.price.toLocaleString()}</td>
                  <td className="potion-rate" style={{ color: rateColor(entry.rate, potion[statKey]) }}>
                    {entry.rate.toFixed(3)}
                  </td>
                  <td className="potion-location">{entry.location}</td>
                </tr>
              ))
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}

export default function PotionTable() {
  return (
    <>
      <PotionSection title="HP 藥水店參考價格" potions={HP_POTIONS} statKey="hp" statLabel="HP" rateColor={hpRateColor} />
      <PotionSection title="MP 藥水店參考價格" potions={MP_POTIONS} statKey="mp" statLabel="MP" rateColor={mpRateColor} />
    </>
  )
}
