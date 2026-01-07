import { useState, useEffect } from 'react';
import '../styles/FilterModal.css';

const FilterModal = ({ isOpen, onClose, onApply, currentFilters }) => {
  const [filters, setFilters] = useState({
    type: currentFilters?.type || '',
    currency: currentFilters?.currency || '',
    priceMin: currentFilters?.priceMin || '',
    priceMax: currentFilters?.priceMax || ''
  });


  useEffect(() => {
    if (isOpen) {
      setFilters({
        type: currentFilters?.type || '',
        currency: currentFilters?.currency || '',
        priceMin: currentFilters?.priceMin || '',
        priceMax: currentFilters?.priceMax || ''
      });
    }
  }, [isOpen, currentFilters]);

  if (!isOpen) return null;

  const handleApply = () => {
    onApply(filters);
    onClose();
  };

  const handleReset = () => {
    const emptyFilters = {
      type: '',
      currency: '',
      priceMin: '',
      priceMax: ''
    };
    setFilters(emptyFilters);
    onApply(emptyFilters);
    onClose();
  };

  return (
    <>
      <div className="filter-modal-overlay" onClick={onClose}>
        <div className="filter-modal" onClick={(e) => e.stopPropagation()}>
          <button className="filter-modal-close" onClick={onClose}>✕</button>

          <h2 className="filter-modal-title">Фильтры</h2>

          <div className="filter-section">
            <label className="filter-label">Тип объявления</label>
            <select
              className="filter-select"
              value={filters.type}
              onChange={(e) => setFilters({ ...filters, type: e.target.value })}
            >
              <option value="">Все</option>
              <option value="Продать">Продать</option>
              <option value="Купить">Купить</option>
              <option value="Сдать в аренду">Сдать в аренду</option>
            </select>
          </div>

          <div className="filter-section">
            <label className="filter-label">Валюта</label>
            <select
              className="filter-select"
              value={filters.currency}
              onChange={(e) => setFilters({ ...filters, currency: e.target.value })}
            >
              <option value="">Все</option>
              <option value="VC">VC</option>
              <option value="$">$</option>
              <option value="BTC">BTC</option>
              <option value="EURO">EURO</option>
            </select>
          </div>

          <div className="filter-section">
            <label className="filter-label">Диапазон цен</label>
            <div className="filter-price-range">
              <input
                type="number"
                className="filter-input"
                placeholder="От"
                value={filters.priceMin}
                onChange={(e) => setFilters({ ...filters, priceMin: e.target.value })}
              />
              <span className="filter-separator">—</span>
              <input
                type="number"
                className="filter-input"
                placeholder="До"
                value={filters.priceMax}
                onChange={(e) => setFilters({ ...filters, priceMax: e.target.value })}
              />
            </div>
          </div>

          <div className="filter-actions">
            <button className="filter-btn filter-btn-apply" onClick={handleApply}>
              Применить
            </button>
            <button className="filter-btn filter-btn-reset" onClick={handleReset}>
              Сбросить
            </button>
          </div>
        </div>
      </div>
    </>
  );
};

export default FilterModal;
